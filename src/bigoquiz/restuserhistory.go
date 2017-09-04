package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"quiz"
	"user"
)

func restHandleUserHistoryAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Use actual authentication.
	var info user.UserHistoryOverall
	info.LoginInfo.Nickname = "example@example.com"

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

	// TODO: Use actual authentication.
	var loginInfo user.LoginInfo
	loginInfo.Nickname = "example@example.com"

	var mapUserStats map[string]*user.UserStats = nil // TODO

	info := buildUserHistorySections(loginInfo, q, mapUserStats)

	jsonStr, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

func buildUserHistorySections(loginInfo user.LoginInfo, quiz *quiz.Quiz, mapUserStats map[string]*user.UserStats) *user.UserHistorySections {
	sections := quiz.Sections
	if sections == nil {
		return nil
	}

	userId := loginInfo.UserId
	quizId := quiz.Id

	var result user.UserHistorySections
	result.LoginInfo = loginInfo
	result.Sections = sections
	result.QuizTitle = quiz.Title

	for _, section := range sections {
		sectionId := section.Id
		if len(sectionId) == 0 {
			continue
		}

		var userStats *user.UserStats = nil
		if mapUserStats != nil {
			userStats = mapUserStats[sectionId]
		}

		if userStats == nil {
			userStats = new(user.UserStats)
			userStats.UserId = userId
			userStats.QuizId = quizId
			userStats.SectionId = sectionId
		}

		fillUserStatsWithTitles(userStats, quiz)
		result.Stats = append(result.Stats, userStats)
	}

	return &result
}

func fillUserStatsWithTitles(userStats *user.UserStats, quiz *quiz.Quiz) {
	sections := quiz.Sections
	if sections == nil {
		return
	}

	// sectionId := userStats.SectionId

	// Set the titles.
	// We don't store these in the datastore because we can get them easily from the Quiz.
	// TODO: for loop over all getTopProblemQuestionHistories() as in gwt-bigoquiz.
}
