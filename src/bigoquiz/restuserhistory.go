package bigoquiz

import (
	"db"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"net/http"
	"quiz"
	"sort"
	"user"
)

// See https://gobyexample.com/sorting-by-functions
type StatsListByTitle []*user.Stats

func (s StatsListByTitle) Len() int {
	return len(s)
}

func (s StatsListByTitle) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s StatsListByTitle) Less(i, j int) bool {
	return s[i].QuizTitle < s[j].QuizTitle
}

func restHandleUserHistoryAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, err := getLoginInfoFromSessionAndDb(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var info user.HistoryOverall
	info.LoginInfo = loginInfo

	// Note: We only show the entire user history for logged-in users,
	// so there is no point in constructing an empty sets of stats for not-logged in users.
	if loginInfo.UserId != nil {
		c := appengine.NewContext(r)
		mapUserStats, err := db.GetUserStats(c, loginInfo.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for quizId, q := range quizzes {
			stats, ok := mapUserStats[quizId]
			if !ok || stats == nil {
				// Show an empty stats section,
				// if there is none in the database yet.
				stats = new(user.Stats)
				stats.QuizId = quizId
			}

			stats.QuizTitle = q.Title
			stats.CountQuestions = q.GetQuestionsCount()
			info.SetQuizStats(q.Id, stats)
		}
	}

	// Sort them alphabetically by quiz title,
	// a a convenience to the client.
	sort.Sort(StatsListByTitle(info.Stats))

	jsonStr, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func restHandleUserHistoryByQuizId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	quizId := ps.ByName(PATH_PARAM_QUIZ_ID)
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	q := getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}

	loginInfo, err := getLoginInfoFromSessionAndDb(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var mapUserStats map[string]*user.Stats
	if loginInfo.UserId != nil {
		c := appengine.NewContext(r)
		mapUserStats, err = db.GetUserStatsForQuiz(c, loginInfo.UserId, quizId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	info := buildUserHistorySections(loginInfo, q, mapUserStats)

	jsonStr, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type Submission struct {
	Answer string `json:"answer"`
}

func restHandleUserHistorySubmitAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string
	var questionId string
	var nextQuestionSectionId string

	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get(QUERY_PARAM_QUIZ_ID)
		questionId = queryValues.Get(QUERY_PARAM_QUESTION_ID)
		nextQuestionSectionId = queryValues.Get(QUERY_PARAM_NEXT_QUESTION_SECTION_ID)
	}

	qa := getQuestionAndAnswer(quizId, questionId)
	if qa == nil {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not parse body.", http.StatusInternalServerError)
		return
	}

	var submission Submission
	err = json.Unmarshal(body, &submission)
	if err != nil {
		http.Error(w, "Could not parse JSON.", http.StatusInternalServerError)
		return
	}

	result := answerIsCorrect(submission.Answer, &qa.Answer)
	submissionResult, err := storeAnswerCorrectnessAndGetSubmissionResult(w, r, quizId, questionId, nextQuestionSectionId, qa, result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.Marshal(submissionResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func restHandleUserHistorySubmitDontKnowAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string
	var questionId string
	var nextQuestionSectionId string

	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get(QUERY_PARAM_QUIZ_ID)
		questionId = queryValues.Get(QUERY_PARAM_QUESTION_ID)
		nextQuestionSectionId = queryValues.Get(QUERY_PARAM_NEXT_QUESTION_SECTION_ID)
	}

	qa := getQuestionAndAnswer(quizId, questionId)
	if qa == nil {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	//Store this like a don't know answer:
	submissionResult, err := storeAnswerCorrectnessAndGetSubmissionResult(w, r, quizId, questionId, nextQuestionSectionId, qa, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.Marshal(submissionResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func restHandleUserHistoryResetSections(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string

	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get(QUERY_PARAM_QUIZ_ID)
	}

	if len(quizId) == 0 {
		http.Error(w, "quiz-id not specified", http.StatusBadRequest)
		return
	}

	q := getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}

	userId, err := getUserIdFromSessionAndDb(r, w)
	if err != nil {
		http.Error(w, "logged-in check failed.", http.StatusInternalServerError)
		return
	}

	if userId == nil {
		loginInfo, err := getLoginInfoFromSessionAndDb(r, w)
		if err != nil {
			http.Error(w, "not logged in. getLoginInfoFromSessionAndDb() returned err.", http.StatusForbidden)
			return
		}

		msg := fmt.Sprintf("not logged in. loginInfo=%v", loginInfo)
		http.Error(w, msg, http.StatusForbidden)
		return
	}

	c := appengine.NewContext(r)
	err = db.DeleteUserStatsForQuiz(c, userId, quizId)
	if err != nil {
		http.Error(w, "deletion of stats failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type SubmissionResult struct {
	Result        bool          `json:"result"`
	CorrectAnswer quiz.Text     `json:"correctAnswer,omitempty"`
	NextQuestion  quiz.Question `json:"nextQuestion,omitempty"`
}

func storeAnswerCorrectnessAndGetSubmissionResult(w http.ResponseWriter, r *http.Request, quizId string, questionId string, nextQuestionSectionId string, qa *quiz.QuestionAndAnswer, result bool) (*SubmissionResult, error) {
	userId, err := getUserIdFromSessionAndDb(r, w)
	if err != nil {
		return nil, fmt.Errorf("getUserIdFromSessionAndDb() failed: %v", err)
	}

	sectionId := qa.Question.SectionId

	// Get the Stats (or a map of them), and use it for both storing the answer and getting the next question,
	// to avoid getting the UserStats twice from the datastore.
	//
	// Call different methods depending on whether nextQuestionSectionId is specified and is the same as the
	// question's section ID, to avoid allocating a map just containing one Stats.
	c := appengine.NewContext(r)
	if nextQuestionSectionId == sectionId {
		var stats *user.Stats
		if userId != nil {
			stats, err = db.GetUserStatsForSection(c, userId, quizId, nextQuestionSectionId)
			if err != nil {
				return nil, fmt.Errorf("GetUserStatsForQuiz() failed: %v", err)
			}

			err := storeAnswerForSection(c, result, quizId, &qa.Question, userId, stats)
			if err != nil {
				return nil, fmt.Errorf("storeAnswerForSection() failed: %v", err)
			}
		}

		return createSubmissionResultForSection(result, quizId, questionId, nextQuestionSectionId, stats)
	} else {
		var stats map[string]*user.Stats
		if userId != nil {
			stats, err = db.GetUserStatsForQuiz(c, userId, quizId)
			if err != nil {
				return nil, fmt.Errorf("GetUserStatsForQuiz() failed: %v", err)
			}

			err := storeAnswer(c, result, quizId, &qa.Question, userId, stats)
			if err != nil {
				return nil, fmt.Errorf("storeAnswer() failed: %v", err)
			}
		}

		return createSubmissionResult(result, quizId, questionId, nextQuestionSectionId, stats)
	}
}

/**
 * stats may be nil.
 */
func createSubmissionResult(result bool, quizId string, questionId string, nextQuestionSectionId string, stats map[string]*user.Stats) (*SubmissionResult, error) {
	q := getQuiz(quizId)

	//We only provide the correct answer if the supplied answer was wrong:
	var correctAnswer *quiz.Text
	if !result {
		correctAnswer = q.GetAnswer(questionId)
	}

	nextQuestion := getNextQuestionFromUserStats(nextQuestionSectionId, q, stats)
	return generateSubmissionResult(result, q, correctAnswer, nextQuestion)
}

func createSubmissionResultForSection(result bool, quizId string, questionId string, nextQuestionSectionId string, stats *user.Stats) (*SubmissionResult, error) {
	q := getQuiz(quizId)

	//We only provide the correct answer if the supplied answer was wrong:
	var correctAnswer *quiz.Text
	if !result {
		correctAnswer = q.GetAnswer(questionId)
	}

	nextQuestion := getNextQuestionFromUserStatsForSection(nextQuestionSectionId, q, stats)
	return generateSubmissionResult(result, q, correctAnswer, nextQuestion)
}

func generateSubmissionResult(result bool, quiz *quiz.Quiz, correctAnswer *quiz.Text, nextQuestion *quiz.Question) (*SubmissionResult, error) {
	var submissionResult SubmissionResult
	submissionResult.Result = result

	if correctAnswer != nil {
		submissionResult.CorrectAnswer = *correctAnswer
	}

	if nextQuestion != nil {
		nextQuestion.SetQuestionExtras(quiz)
		submissionResult.NextQuestion = *nextQuestion
	}

	return &submissionResult, nil
}

func getNextQuestionFromUserStats(sectionId string, q *quiz.Quiz, stats map[string]*user.Stats) *quiz.Question {
	const MAX_TRIES int = 10
	var tries int
	var question *quiz.Question
	var questionBestSoFar *quiz.Question
	var questionBestCountAnsweredWrong int

	for tries < MAX_TRIES {
		tries += 1

		question = q.GetRandomQuestion(sectionId)
		if question == nil {
			continue
		}

		if questionBestSoFar == nil {
			questionBestSoFar = question
		}

		if stats == nil {
			//Assume this means the user has never answered any question in any section.
			return question
		}

		userStats, ok := stats[question.SectionId]
		if !ok || userStats == nil {
			//Assume this means the user has never answered any question in the section.
			return question
		}

		questionId := question.Id

		//Prioritize questions that have never been asked.
		if !userStats.GetQuestionWasAnswered(questionId) {
			return question
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

	return questionBestSoFar
}

/** stats may be nil
 */
func getNextQuestionFromUserStatsForSection(sectionId string, quiz *quiz.Quiz, stats *user.Stats) *quiz.Question {
	//TODO: Avoid this temporary map:
	m := make(map[string]*user.Stats)

	if stats != nil {
		m[stats.SectionId] = stats
	}

	return getNextQuestionFromUserStats(sectionId, quiz, m)
}

/** Update the user.Stats for the question's quiz section, in the database,
 * storing a new user.Stats in the database if necessary.
 * To update an existing user.Stats in the database, via the stats parameter, its Key field must be set.
 */
func storeAnswer(c context.Context, result bool, quizId string, question *quiz.Question, userId *datastore.Key, stats map[string]*user.Stats) error {
	if userId == nil {
		return fmt.Errorf("storeAnswer(): userId is nil")
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
		sectionStats = new(user.Stats)
		sectionStats.UserId = userId
		sectionStats.QuizId = quizId
		sectionStats.SectionId = sectionId
	}

	return storeAnswerForSection(c, result, quizId, question, userId, sectionStats)
}

/** Update the user.Stats for the section, for the quiz, in the database,
 * storing a new user.Stats in the database if necessary.
 * To update an existing user.Stats in the database, its Key field must be set.
 */
func storeAnswerForSection(c context.Context, result bool, quizId string, question *quiz.Question, userId *datastore.Key, sectionStats *user.Stats) error {
	if userId == nil {
		return fmt.Errorf("storeAnswerForSection(): userId is nil")
	}

	if question == nil {
		return fmt.Errorf("storeAnswerForSection(): question is nil")
	}

	sectionId := question.SectionId
	if len(sectionId) == 0 {
		return fmt.Errorf("storeAnswerForSection(): question's section ID is empty")
	}

	if sectionStats == nil {
		sectionStats = new(user.Stats)
		sectionStats.UserId = userId
		sectionStats.QuizId = quizId
		sectionStats.SectionId = sectionId
	}

	sectionStats.IncrementAnswered()

	if result {
		sectionStats.IncrementCorrect()
	}

	sectionStats.UpdateProblemQuestion(question, result)

	if err := db.StoreUserStats(c, sectionStats); err != nil {
		return fmt.Errorf("db.StoreUserStat() for key: %v: %v", sectionStats.Key, err)
	}

	return nil
}

func answerIsCorrect(answer string, correctAnswer *quiz.Text) bool {
	if correctAnswer == nil {
		return false
	}

	return correctAnswer.Text == answer
}

func getQuestionAndAnswer(quizId string, questionId string) *quiz.QuestionAndAnswer {

	q := getQuiz(quizId)
	if q == nil {
		return nil
	}

	return q.GetQuestionAndAnswer(questionId)
}

func buildUserHistorySections(loginInfo *user.LoginInfo, quiz *quiz.Quiz, mapUserStats map[string]*user.Stats) *user.HistorySections {
	sections := quiz.Sections
	if sections == nil {
		return nil
	}

	userId := loginInfo.UserId
	quizId := quiz.Id

	var result user.HistorySections
	result.LoginInfo = *loginInfo
	result.QuizId = quizId
	result.QuizTitle = quiz.Title

	for _, section := range sections {
		sectionId := section.Id
		if len(sectionId) == 0 {
			continue
		}

		var userStats *user.Stats = nil
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
			userStats = new(user.Stats)
			userStats.UserId = userId
			userStats.QuizId = quizId
			userStats.SectionId = sectionId
		}

		userStats.QuizTitle = quiz.Title
		userStats.SectionTitle = section.Title
		userStats.CountQuestions = section.CountQuestions

		fillUserStatsWithExtras(userStats, quiz)
		result.Stats = append(result.Stats, userStats)
	}

	return &result
}

func fillUserStatsWithExtras(userStats *user.Stats, qz *quiz.Quiz) {
	// Set the titles.
	// We don't store these in the datastore because we can get them easily from the Quiz.

	// TODO: Only send the top problem question histories in the JSON,
	// instead of all of them?
	for i := range userStats.QuestionHistories {
		qh := &(userStats.QuestionHistories[i])

		q := qz.GetQuestionAndAnswer(qh.QuestionId)
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
}
