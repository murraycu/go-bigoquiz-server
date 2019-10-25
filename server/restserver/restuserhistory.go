package restserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	domainuser "github.com/murraycu/go-bigoquiz-server/domain/user"
	restquiz "github.com/murraycu/go-bigoquiz-server/server/restserver/quiz"
	restuser "github.com/murraycu/go-bigoquiz-server/server/restserver/user"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

// See https://gobyexample.com/sorting-by-functions
type StatsListByTitle []*restuser.Stats

func (s StatsListByTitle) Len() int {
	return len(s)
}

func (s StatsListByTitle) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s StatsListByTitle) Less(i, j int) bool {
	return s[i].QuizTitle < s[j].QuizTitle
}

func (s *RestServer) HandleUserHistoryAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, userId, err := s.getLoginInfoFromSessionAndDb(r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "getLoginInfoFromSessionAndDb() failed: %v", err)
		return
	}

	var info restuser.HistoryOverall
	info.LoginInfo = loginInfo

	// Note: We only show the entire user history for logged-in users,
	// so there is no point in constructing an empty sets of stats for not-logged in users.
	if loginInfo.LoggedIn && len(userId) != 0 {
		c := r.Context()

		mapUserStats, err := s.userDataClient.GetUserStats(c, userId)
		if err != nil {
			handleErrorAsHttpError(w, http.StatusInternalServerError, "GetUserStats() failed: %v", err)
			return
		}

		for _, q := range s.quizzesListSimple {
			quizId := q.Id
			stats, ok := mapUserStats[quizId]
			if !ok || stats == nil {
				// Show an empty stats section,
				// if there is none in the database yet.
				stats = new(domainuser.Stats)
				stats.QuizId = quizId
			}

			quizCache, ok := s.quizCacheMap[q.Id]
			if !ok {
				continue
			}

			restStats, err := convertDomainStatsToRestStats(stats, quizCache)
			if err != nil {
				handleErrorAsHttpError(w, http.StatusInternalServerError, "convertDomainStatsToRestStats() failed: %v", err)
				return
			}

			// Set extras for REST:
			restStats.QuizTitle = q.Title
			restStats.CountQuestions = quizCache.GetQuestionsCount()

			info.SetQuizStats(q.Id, restStats)
		}
	}

	// Sort them alphabetically by quiz title,
	// a a convenience to the client.
	sort.Sort(StatsListByTitle(info.Stats))

	marshalAndWriteOrHttpError(w, &info)
}

func (s *RestServer) HandleUserHistoryByQuizId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	quizId := ps.ByName(PATH_PARAM_QUIZ_ID)
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		handleErrorAsHttpError(w, http.StatusInternalServerError, "Empty quiz ID")
		return
	}

	q := s.getQuiz(quizId)
	if q == nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "quiz not found")
		return
	}

	loginInfo, userId, err := s.getLoginInfoFromSessionAndDb(r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "getLoginInfoFromSessionAndDb() failed: %v", err)
		return
	}

	var mapUserStats map[string]*domainuser.Stats
	if loginInfo.LoggedIn && len(userId) != 0 {
		c := r.Context()

		mapUserStats, err = s.userDataClient.GetUserStatsForQuiz(c, userId, quizId)
		if err != nil {
			handleErrorAsHttpError(w, http.StatusInternalServerError, "GetUserStatsForQuiz() falied: %v", err)
			return
		}
	}

	info, err := s.buildRestUserHistorySections(loginInfo, q, mapUserStats)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "buildRestUserHistorySections() failed: %v", err)
		return
	}

	marshalAndWriteOrHttpError(w, &info)
}

type Submission struct {
	Answer string `json:"answer"`
}

func (s *RestServer) HandleUserHistorySubmitAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string
	var questionId string
	var nextQuestionSectionId string

	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get(QUERY_PARAM_QUIZ_ID)
		questionId = queryValues.Get(QUERY_PARAM_QUESTION_ID)
		nextQuestionSectionId = queryValues.Get(QUERY_PARAM_NEXT_QUESTION_SECTION_ID)
	}

	qa, err := s.getQuestionAndAnswer(quizId, questionId)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "question not found")
		return
	}

	if qa == nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "question not found")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "Could not parse body. ioutil.ReadAll() failed: %v", err)
		return
	}

	var submission Submission
	err = json.Unmarshal(body, &submission)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "Could not parse JSON. json.Unmarshal() failed: %v", err)
		return
	}

	userId, err := s.getUserIdFromSessionAndDb(r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "getUserIdFromSessionAndDb() failed: %v", err)
	}

	result := answerIsCorrect(submission.Answer, &qa.Answer)
	submissionResult, err := s.storeAnswerCorrectnessAndGetSubmissionResult(r.Context(), userId, quizId, nextQuestionSectionId, qa, result)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "storeAnswerCorrectnessAndGetSubmissionResult() failed: %v", err)
		return
	}

	marshalAndWriteOrHttpError(w, &submissionResult)
}

func (s *RestServer) HandleUserHistorySubmitDontKnowAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string
	var questionId string
	var nextQuestionSectionId string

	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get(QUERY_PARAM_QUIZ_ID)
		questionId = queryValues.Get(QUERY_PARAM_QUESTION_ID)
		nextQuestionSectionId = queryValues.Get(QUERY_PARAM_NEXT_QUESTION_SECTION_ID)
	}

	qa, err := s.getQuestionAndAnswer(quizId, questionId)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "question not found")
		return
	}

	if qa == nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "question not found")
		return
	}

	userId, err := s.getUserIdFromSessionAndDb(r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "getUserIdFromSessionAndDb() failed: %v", err)
	}

	//Store this like a don't know answer:
	submissionResult, err := s.storeAnswerCorrectnessAndGetSubmissionResult(r.Context(), userId, quizId, nextQuestionSectionId, qa, false)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "storeAnswerCorrectnessAndGetSubmissionResult() failed: %v", err)
		return
	}

	marshalAndWriteOrHttpError(w, &submissionResult)
}

func (s *RestServer) HandleUserHistoryResetSections(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string

	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get(QUERY_PARAM_QUIZ_ID)
	}

	if len(quizId) == 0 {
		handleErrorAsHttpError(w, http.StatusBadRequest, "quiz-id not specified")
		return
	}

	q := s.getQuiz(quizId)
	if q == nil {
		handleErrorAsHttpError(w, http.StatusNotFound, "quiz not found")
		return
	}

	userId, err := s.getUserIdFromSessionAndDb(r)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "logged-in check failed. getUserIdFromSessionAndDb() failed: %v", err)
		return
	}

	if len(userId) == 0 {
		loginInfo, _, err := s.getLoginInfoFromSessionAndDb(r)
		if err != nil {
			handleErrorAsHttpError(w, http.StatusForbidden, "not logged in. getLoginInfoFromSessionAndDb() failed: %v", err)
			return
		}

		msg := fmt.Sprintf("not logged in. loginInfo=%v", loginInfo)
		handleErrorAsHttpError(w, http.StatusForbidden, msg)
		return
	}

	c := r.Context()

	err = s.userDataClient.DeleteUserStatsForQuiz(c, userId, quizId)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "deletion of stats failed: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type SubmissionResult struct {
	Result        bool              `json:"result"`
	CorrectAnswer restquiz.Text     `json:"correctAnswer,omitempty"`
	NextQuestion  restquiz.Question `json:"nextQuestion,omitempty"`
}

func (s *RestServer) storeAnswerCorrectnessAndGetSubmissionResult(c context.Context, userId string, quizId string, nextQuestionSectionId string, qa *restquiz.QuestionAndAnswer, result bool) (*SubmissionResult, error) {
	sectionId := qa.Question.SectionId
	questionId := qa.Question.Id

	// Get the Stats (or a map of them), and use it for both storing the answer and getting the next question,
	// to avoid getting the UserStats twice from the datastore.
	//
	// Call different methods depending on whether nextQuestionSectionId is specified and is the same as the
	// question's section ID, to avoid allocating a map just containing one Stats.
	if nextQuestionSectionId == sectionId {
		var stats *domainuser.Stats
		if len(userId) != 0 {
			var err error
			stats, err = s.userDataClient.GetUserStatsForSection(c, userId, quizId, nextQuestionSectionId)
			if err != nil {
				return nil, fmt.Errorf("GetUserStatsForQuiz() failed: %v", err)
			}

			err = s.storeAnswerForSection(c, result, quizId, &qa.Question, userId, stats)
			if err != nil {
				return nil, fmt.Errorf("storeAnswerForSection() failed: %v", err)
			}
		}

		result, err := s.createSubmissionResultForSection(result, quizId, questionId, nextQuestionSectionId, stats)
		if err != nil {
			return nil, fmt.Errorf("createSubmissionResultForSection() failed: %v", err)
		}

		return result, nil
	} else {
		var stats map[string]*domainuser.Stats
		if len(userId) != 0 {
			stats, err := s.userDataClient.GetUserStatsForQuiz(c, userId, quizId)
			if err != nil {
				return nil, fmt.Errorf("GetUserStatsForQuiz() failed: %v", err)
			}

			err = s.storeAnswer(c, result, quizId, &qa.Question, userId, stats)
			if err != nil {
				return nil, fmt.Errorf("storeAnswer() failed: %v", err)
			}
		}

		result, err := s.createSubmissionResult(result, quizId, questionId, nextQuestionSectionId, stats)
		if err != nil {
			return nil, fmt.Errorf("createSubmissionResultForSection() failed: %v", err)
		}

		return result, nil
	}
}

/**
 * stats may be nil.
 */
func (s *RestServer) createSubmissionResult(result bool, quizId string, questionId string, nextQuestionSectionId string, stats map[string]*domainuser.Stats) (*SubmissionResult, error) {
	quizCache, err := s.getQuizCache(quizId)
	if err != nil {
		return nil, fmt.Errorf("couldn't find quiz in quizCacheMap with quiz ID: %v: %v", quizId, err)
	}

	//We only provide the correct answer if the supplied answer was wrong:
	var correctAnswer *restquiz.Text
	if !result {
		correctAnswer = quizCache.GetAnswer(questionId)
	}

	q := s.getQuiz(quizId)
	if q == nil {
		return nil, fmt.Errorf("Couldn't find quiz with quiz ID: %v", quizId)
	}

	nextQuestion, err := s.getNextQuestionFromUserStats(nextQuestionSectionId, q, stats)
	if err != nil {
		return nil, fmt.Errorf("getNextQuestionFromUserStats() failed: %v", err)
	}

	return s.generateSubmissionResult(result, quizCache, correctAnswer, nextQuestion)
}

func (s *RestServer) createSubmissionResultForSection(result bool, quizId string, questionId string, nextQuestionSectionId string, stats *domainuser.Stats) (*SubmissionResult, error) {
	quizCache, err := s.getQuizCache(quizId)
	if err != nil {
		return nil, fmt.Errorf("couldn't find quiz in quizCacheMap with quiz ID: %v: %v", quizId, err)
	}

	//We only provide the correct answer if the supplied answer was wrong:
	var correctAnswer *restquiz.Text
	if !result {
		correctAnswer = quizCache.GetAnswer(questionId)
	}

	q := s.getQuiz(quizId)
	if q == nil {
		return nil, fmt.Errorf("Couldn't find quiz with quiz ID: %v", quizId)
	}

	nextQuestion, err := s.getNextQuestionFromUserStatsForSection(nextQuestionSectionId, q, stats)
	if err != nil {
		return nil, fmt.Errorf("getNextQuestionFromUserStatsForSection() failed: %v", err)
	}

	submissionresult, err := s.generateSubmissionResult(result, quizCache, correctAnswer, nextQuestion)
	if err != nil {
		return nil, fmt.Errorf("generateSubmissionResult() failed: %v", err)
	}

	return submissionresult, nil
}

func (s *RestServer) generateSubmissionResult(result bool, quizCache *QuizCache, correctAnswer *restquiz.Text, nextQuestion *restquiz.Question) (*SubmissionResult, error) {
	var submissionResult SubmissionResult
	submissionResult.Result = result

	if correctAnswer != nil {
		submissionResult.CorrectAnswer = *correctAnswer
	}

	if nextQuestion != nil {
		submissionResult.NextQuestion = *nextQuestion
	}

	return &submissionResult, nil
}

func (s *RestServer) getNextQuestionFromUserStats(sectionId string, q *restquiz.Quiz, stats map[string]*domainuser.Stats) (*restquiz.Question, error) {
	const MAX_TRIES int = 10
	var tries int
	var question *restquiz.Question
	var questionBestSoFar *restquiz.Question
	var questionBestCountAnsweredWrong int

	for tries < MAX_TRIES {
		tries += 1

		var err error
		question, err = s.GetRandomQuestion(q, sectionId)
		if err != nil {
			return nil, fmt.Errorf("GetRandomQuestion() failed: %v", err)
		}

		if question == nil {
			continue
		}

		if questionBestSoFar == nil {
			questionBestSoFar = question
		}

		if stats == nil {
			//Assume this means the user has never answered any question in any section.
			return question, nil
		}

		userStats, ok := stats[question.SectionId]
		if !ok || userStats == nil {
			//Assume this means the user has never answered any question in the section.
			return question, nil
		}

		questionId := question.Id

		//Prioritize questions that have never been asked.
		if !userStats.GetQuestionWasAnswered(questionId) {
			return question, nil
		}

		//Otherwise, try a few times to get a question that
		//we have got wrong many times:
		//We could just get the most-wrong answer directly,
		//but we want some randomness.
		countAnsweredWrong := userStats.GetQuestionCountAnsweredWrong(questionId)
		if countAnsweredWrong > questionBestCountAnsweredWrong {
			questionBestSoFar = question
			questionBestCountAnsweredWrong = countAnsweredWrong
		}
	}

	return questionBestSoFar, nil
}

/** stats may be nil
 */
func (s *RestServer) getNextQuestionFromUserStatsForSection(
	sectionId string, quiz *restquiz.Quiz, stats *domainuser.Stats) (*restquiz.Question, error) {
	//TODO: Avoid this temporary map:
	m := make(map[string]*domainuser.Stats)

	if stats != nil {
		m[stats.SectionId] = stats
	}

	result, err := s.getNextQuestionFromUserStats(sectionId, quiz, m)
	if err != nil {
		return nil, fmt.Errorf("getNextQuestionFromUserStats() failed: %v", err)
	}

	return result, nil
}

/** Update the user.Stats for the question's quiz section, in the database,
 * storing a new user.Stats in the database if necessary.
 * To update an existing user.Stats in the database, via the stats parameter, its Key field must be set.
 */
func (s *RestServer) storeAnswer(c context.Context, result bool, quizId string, question *restquiz.Question, userId string, stats map[string]*domainuser.Stats) error {
	if len(userId) == 0 {
		return fmt.Errorf("storeAnswer(): userId is empty")
	}

	if question == nil {
		return fmt.Errorf("storeAnswer(): question is nil")
	}

	sectionId := question.SectionId
	if len(sectionId) == 0 {
		return fmt.Errorf("storeAnswer(): question's section ID is empty")
	}

	sectionStats, ok := stats[sectionId]
	if !ok {
		// It's not in the map yet, so we add it.
		sectionStats = new(domainuser.Stats)
		sectionStats.QuizId = quizId
		sectionStats.SectionId = sectionId
	}

	return s.storeAnswerForSection(c, result, quizId, question, userId, sectionStats)
}

/** Update the user.Stats for the section, for the quiz, in the database,
 * storing a new user.Stats in the database if necessary.
 * To update an existing user.Stats in the database, its Key field must be set.
 */
func (s *RestServer) storeAnswerForSection(c context.Context, result bool, quizId string, question *restquiz.Question, userId string, sectionStats *domainuser.Stats) error {
	if len(userId) == 0 {
		return fmt.Errorf("storeAnswerForSection(): userId is empty")
	}

	if question == nil {
		return fmt.Errorf("storeAnswerForSection(): question is nil")
	}

	sectionId := question.SectionId
	if len(sectionId) == 0 {
		return fmt.Errorf("storeAnswerForSection(): question's section ID is empty")
	}

	if sectionStats == nil {
		sectionStats = new(domainuser.Stats)
		sectionStats.QuizId = quizId
		sectionStats.SectionId = sectionId
	}

	sectionStats.UpdateStatsForAnswerCorrectness(question.Id, result)

	if err := s.userDataClient.StoreUserStats(c, userId, sectionStats); err != nil {
		return fmt.Errorf("db.StoreUserStat() failed for: %v: %v", sectionStats, err)
	}

	return nil
}

func answerIsCorrect(answer string, correctAnswer *restquiz.Text) bool {
	if correctAnswer == nil {
		return false
	}

	return correctAnswer.Text == answer
}

func (s *RestServer) getQuestionAndAnswer(quizId string, questionId string) (*restquiz.QuestionAndAnswer, error) {
	quizCache, err := s.getQuizCache(quizId)
	if err != nil {
		return nil, fmt.Errorf("could not find quiz cache for quiz with ID: %v: %v", quizId, err)
	}

	return quizCache.GetQuestionAndAnswer(questionId), nil
}

func (s *RestServer) buildRestUserHistorySections(loginInfo *restuser.LoginInfo, quiz *restquiz.Quiz, mapUserStats map[string]*domainuser.Stats) (*restuser.HistorySections, error) {
	sections := quiz.Sections
	if sections == nil {
		return nil, nil
	}

	quizId := quiz.Id

	var result restuser.HistorySections
	result.LoginInfo = *loginInfo
	result.QuizId = quizId
	result.QuizTitle = quiz.Title

	for _, section := range sections {
		sectionId := section.Id
		if len(sectionId) == 0 {
			continue
		}

		var userStats *domainuser.Stats = nil
		if mapUserStats != nil {
			var ok bool
			userStats, ok = mapUserStats[sectionId]

			// A sanity check:
			if ok && userStats != nil && userStats.SectionId != sectionId {
				// This does not happen, but let's be sure.
				userStats = nil
			}
		}

		if userStats == nil {
			userStats = new(domainuser.Stats)
			userStats.QuizId = quizId
			userStats.SectionId = sectionId
		}

		quizCache, err := s.getQuizCache(quizId)
		if err != nil {
			return nil, fmt.Errorf("could not find quiz cache for quiz with ID: %v", quizId)
		}

		restStats, err := convertDomainStatsToRestStats(userStats, quizCache)
		if err != nil {
			// Log this, but forgive it. Maybe the quiz has changed.
			log.Printf("convertDomainStatsToRestStats() failed: %v", err)
			continue
		}

		// Set extras for REST:
		restStats.QuizTitle = quiz.Title
		restStats.SectionTitle = section.Title
		restStats.CountQuestions = quizCache.GetQuestionsCount()

		err = s.fillUserStatsWithExtras(restStats, quiz)
		if err != nil {
			return nil, fmt.Errorf("fillUserStatsWithExtras() failed: %v", err)
		}

		result.Stats = append(result.Stats, restStats)
	}

	return &result, nil
}

func (s *RestServer) fillUserStatsWithExtras(userStats *restuser.Stats, qz *restquiz.Quiz) error {
	// Set the titles.
	// We don't store these in the datastore because we can get them easily from the Quiz.

	quizCache, err := s.getQuizCache(qz.Id)
	if err != nil {
		return fmt.Errorf("couldn't find quiz in quizCacheMap with quiz ID: %v: %v", qz.Id, err)
	}

	// TODO: Only send the top problem question histories in the JSON,
	// instead of all of them?
	for i := range userStats.QuestionHistories {
		qh := &(userStats.QuestionHistories[i])

		q := quizCache.GetQuestionAndAnswer(qh.QuestionId)
		if q == nil {
			continue
		}

		// Extras that are useful to the client via JSON
		// but not stored in the database.
		qh.QuestionTitle = &(q.Text)
		qh.SectionId = q.SectionId

		if q.SubSection != nil {
			qh.SubSectionTitle = q.SubSection.Title
		} else {
			qh.SubSectionTitle = ""
		}
	}

	return nil
}
