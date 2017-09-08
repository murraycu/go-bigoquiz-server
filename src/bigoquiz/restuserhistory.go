package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"quiz"
	"user"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"db"
	"google.golang.org/appengine"
)

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
		info.LoginInfo.ErrorMessage = "debug2 message"

		c := appengine.NewContext(r)
		mapUserStats, err := db.GetUserStats(c, loginInfo.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for quizId, q := range quizzes {
			stats, ok := mapUserStats[quizId]
			if (!ok || stats == nil) {
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

	jsonStr, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

func restHandleUserHistoryByQuizId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	quizId := ps.ByName("quizId")
	if quizId == "" {
		// This makes no sense. restHandleQuizAll() should have been called.
		http.Error(w, "Empty quiz ID", http.StatusInternalServerError)
		return
	}

	q := getQuiz(quizId)
	if q == nil {
		http.Error(w, "quiz not found", http.StatusInternalServerError)
		return
	}

	loginInfo, err := getLoginInfoFromSessionAndDb(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var mapUserStats map[string]*user.Stats;
	if loginInfo.UserId == nil {
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

	w.Write(jsonStr)
}

func restHandleUserHistorySubmitAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string
	var questionId string
	var answer string
	var nextQuestionSectionId string

	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get("quiz-id")
		questionId = queryValues.Get("question-id")
		answer = queryValues.Get("answer")
		nextQuestionSectionId = queryValues.Get("next-question-section-id")
	}

	qa := getQuestionAndAnswer(quizId, questionId)
	if qa == nil {
		http.Error(w, "question not found", http.StatusInternalServerError)
		return
	}

	result := answerIsCorrect(answer, &qa.Answer)
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

	w.Write(jsonStr)
}

func restHandleUserHistorySubmitDontKnowAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var quizId string
	var questionId string
	var nextQuestionSectionId string

	queryValues := r.URL.Query()
	if queryValues != nil {
		quizId = queryValues.Get("quiz-id")
		questionId = queryValues.Get("question-id")
		nextQuestionSectionId = queryValues.Get("next-question-section-id")
	}

	qa := getQuestionAndAnswer(quizId, questionId)
	if qa == nil {
		http.Error(w, "question not found", http.StatusInternalServerError)
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

	w.Write(jsonStr)
}


type SubmissionResult struct {
	Result bool `json:"result"`
	CorrectAnswer quiz.Text `json:"correctAnswer,omitempty"`
	NextQuestion quiz.Question `json:"nextQuestion,omitempty"`
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
		stats, err := db.GetUserStatsForSection(c, userId, quizId, nextQuestionSectionId)
		if err != nil {
			return nil, fmt.Errorf("GetUserStatsForQuiz() failed: %v", err)
		}

		storeAnswerForSection(c, result, quizId, &qa.Question, userId, stats)

		return createSubmissionResultForSection(result, quizId, questionId, nextQuestionSectionId, stats)
	} else {
		stats, err := db.GetUserStatsForQuiz(c, userId, quizId)
		if err != nil {
			return nil, fmt.Errorf("GetUserStatsForQuiz() failed: %v", err)
		}

		storeAnswer(c, result, quizId, &qa.Question, userId, stats)

		return createSubmissionResult(result, quizId, questionId, nextQuestionSectionId, stats)
	}
}

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
		setQuestionExtras(nextQuestion, quiz)
		submissionResult.NextQuestion = *nextQuestion
	}

	return &submissionResult, nil
}

// Set extra details for the question,
// such as quiz and section titles,
// so the client doesn't need to look these up separately.
func setQuestionExtras(question *quiz.Question, q *quiz.Quiz) {
	var briefSection *quiz.HasIdAndTitle
	var subSection *quiz.SubSection

	section := q.GetSection(question.SectionId)
	if section != nil {
		// Create a simpler version of the Section information:
		briefSection = new(quiz.HasIdAndTitle)
		briefSection.Id = section.Id
		briefSection.Title = section.Title
		briefSection.Link = section.Link

		subSection = section.GetSubSection(question.SubSectionId)
	}

	question.SetTitles(q.Title, briefSection, subSection)

	question.QuizUsesMathML = q.UsesMathML
}

func getNextQuestionFromUserStats(sectionId string, quiz *quiz.Quiz, stats map[string]*user.Stats) *quiz.Question {
	return quiz.GetRandomQuestion() //TODO
}

func getNextQuestionFromUserStatsForSection(sectionId string, quiz *quiz.Quiz, stats *user.Stats) *quiz.Question {
	return quiz.GetRandomQuestion() //TODO
}

func storeAnswer(c context.Context, result bool, quizId string, question *quiz.Question, userId *datastore.Key, stats map[string]*user.Stats) error {
	if userId == nil {
		return fmt.Errorf("storeAnswer(): userId is nil.")
	}

	if question == nil {
		return fmt.Errorf("storeAnswer(): question is nil.")
	}

	sectionId := question.SectionId
	if len(sectionId) == 0 {
		return fmt.Errorf("storeAnswer(): question's section ID is empty.")
	}

	sectionStats, ok := stats[sectionId]
	if !ok {
		// It's not in the map yet, so we add it.
		sectionStats = new(user.Stats)
		sectionStats.UserId = userId
		sectionStats.QuizId = quizId
		sectionStats.SectionId = sectionId
	}

	storeAnswerForSection(c, result, quizId, question, userId, sectionStats)
	return nil
}

func storeAnswerForSection(c context.Context, result bool, quizId string, question *quiz.Question, userId *datastore.Key, sectionStats *user.Stats) error {
	if userId == nil {
		return fmt.Errorf("storeAnswerForSection(): userId is nil.")
	}

	if question == nil {
		return fmt.Errorf("storeAnswerForSection(): question is nil.")
	}

	sectionId := question.SectionId
	if len(sectionId) == 0 {
		return fmt.Errorf("storeAnswerForSection(): question's section ID is empty.")
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
		return fmt.Errorf("db.StoreUserStat(): %v", sectionStats.Key, err)
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
	result.QuizTitle = quiz.Title

	for _, section := range sections {
		sectionId := section.Id
		if len(sectionId) == 0 {
			continue
		}

		var userStats *user.Stats = nil
		if mapUserStats != nil {
			userStats = mapUserStats[sectionId]
		}

		if userStats == nil {
			userStats = new(user.Stats)
			userStats.UserId = userId
			userStats.QuizId = quizId
			userStats.SectionId = sectionId

			userStats.QuizTitle = quiz.Title
			userStats.SectionTitle = section.Title
			// TODO: userStats.CountQuestions = section.GetQuestionsCount()
		}

		fillUserStatsWithTitles(userStats, quiz)
		result.Stats = append(result.Stats, userStats)
	}

	return &result
}

func fillUserStatsWithTitles(userStats *user.Stats, quiz *quiz.Quiz) {
	sections := quiz.Sections
	if sections == nil {
		return
	}

	// sectionId := userStats.SectionId

	// Set the titles.
	// We don't store these in the datastore because we can get them easily from the Quiz.
	// TODO: for loop over all getTopProblemQuestionHistories() as in gwt-bigoquiz.
}
