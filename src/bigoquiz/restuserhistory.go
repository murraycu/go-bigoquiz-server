package bigoquiz

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"quiz"
	"user"
)

func restHandleUserHistoryAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginInfo, err := getLoginInfoFromSessionAndDb(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var info user.HistoryOverall
	info.LoginInfo = *loginInfo

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

	var mapUserStats map[string]*user.Stats = nil // TODO

	info := buildUserHistorySections(loginInfo, q, mapUserStats)

	jsonStr, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
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
	result.Sections = sections
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
