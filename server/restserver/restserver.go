package restserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/murraycu/go-bigoquiz-server/repositories/db"
	"github.com/murraycu/go-bigoquiz-server/repositories/quizzes"
	restquiz "github.com/murraycu/go-bigoquiz-server/server/restserver/quiz"
	"github.com/murraycu/go-bigoquiz-server/server/usersessionstore"
)

const QUERY_PARAM_QUIZ_ID = "quiz-id"
const QUERY_PARAM_SECTION_ID = "section-id"
const QUERY_PARAM_QUESTION_ID = "question-id"
const QUERY_PARAM_LIST_ONLY = "list-only"
const QUERY_PARAM_NEXT_QUESTION_SECTION_ID = "next-question-section-id"
const PATH_PARAM_QUIZ_ID = "quizId"
const PATH_PARAM_QUESTION_ID = "questionId"

type restQuizList []*restquiz.Quiz

// A map of quiz ID to the quiz's restQuestionsMap
type restQuizCacheMap map[string]*QuizCache

type RestServer struct {
	quizzes           restQuizMap
	quizzesListSimple restQuizList
	quizzesListFull   restQuizList

	// Easier access to some quiz details.
	quizCacheMap restQuizCacheMap

	userDataClient db.UserDataRepository

	// Session cookie store.
	userSessionStore usersessionstore.UserSessionStore
}

func NewRestServer(quizzesStore quizzes.QuizzesRepository, userSessionStore usersessionstore.UserSessionStore, userDataRepository db.UserDataRepository) (*RestServer, error) {
	result := &RestServer{}
	result.userDataClient = userDataRepository

	quizzes, err := quizzesStore.LoadQuizzes()
	if err != nil {
		return nil, fmt.Errorf("LoadQuizzes() failed: %v", err)
	}

	result.quizzes, err = convertDomainQuizzesToRestQuizzes(quizzes)
	if err != nil {
		return nil, fmt.Errorf("convertDomainQuizzesToRestQuizzes() failed: %v", err)
	}

	// Fill the QuizCache map.
	result.quizCacheMap = make(restQuizCacheMap)
	for _, q := range result.quizzes {
		quizCache, err := NewQuizCache(q)
		if err != nil {
			return nil, fmt.Errorf("NewQuizCache() failed: %v", err)
		}

		result.quizCacheMap[q.Id] = quizCache

		// Use it to fill the extras:
		err = fillRestQuizExtrasFromQuizCache(q, quizCache)
		if err != nil {
			return nil, fmt.Errorf("fillRestQuizzesExtrasFromQuizCache() failed: %v", err)
		}
	}

	result.quizzesListSimple = buildQuizzesSimple(result.quizzes)
	result.quizzesListFull = buildQuizzesFull(result.quizzes)

	result.userSessionStore = userSessionStore

	return result, nil
}

// See https://gobyexample.com/sorting-by-functions
type quizListByTitle []*restquiz.Quiz

func (s quizListByTitle) Len() int {
	return len(s)
}

func (s quizListByTitle) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s quizListByTitle) Less(i, j int) bool {
	return s[i].Title < s[j].Title
}

// TODO: Is there instead some way to output just the top-level of the JSON,
// and only some of the fields?
func buildQuizzesSimple(quizzes restQuizMap) restQuizList {
	// Create a slice with the same capacity.
	result := make(restQuizList, 0, len(quizzes))

	for _, q := range quizzes {
		var simple restquiz.Quiz

		simple.Id = q.Id
		simple.Title = q.Title
		simple.Link = q.Link

		simple.IsPrivate = q.IsPrivate

		result = append(result, &simple)
	}

	sort.Sort(quizListByTitle(result))

	return result
}

func buildQuizzesFull(quizzes restQuizMap) restQuizList {
	// Create a slice with the same capacity.
	result := make(restQuizList, 0, len(quizzes))

	for _, q := range quizzes {
		result = append(result, q)
	}

	sort.Sort(quizListByTitle(result))

	return result
}

func (self *RestServer) GetRandomQuestion(quiz *restquiz.Quiz, sectionId string) (*restquiz.Question, error) {
	quizCache, err := self.getQuizCache(quiz.Id)
	if err != nil {
		return nil, fmt.Errorf("getQuizCache() failed: %v", err)
	}

	return quizCache.GetRandomQuestion(sectionId), nil
}

func (self *RestServer) getQuizCache(quizId string) (*QuizCache, error) {
	if self.quizCacheMap == nil {
		return nil, fmt.Errorf("quizCacheMap is nil")
	}

	quizCache, ok := self.quizCacheMap[quizId]
	if !ok {
		return nil, fmt.Errorf("could not find Quiz cache for quiz ID: %v", quizId)
	}

	return quizCache, nil
}

/** Set extra details for the question,
 * such as quiz and section titles,
 * so the client doesn't need to look these up separately.
 */
func setQuestionExtras(question *restquiz.Question, section *restquiz.Section, subSection *restquiz.SubSection, quizCache *QuizCache) error {
	var briefSection *restquiz.HasIdAndTitle

	if section != nil {
		// Create a simpler version of the Section information:
		briefSection = new(restquiz.HasIdAndTitle)
		briefSection.Id = section.Id
		briefSection.Title = section.Title
		briefSection.Link = section.Link
	}

	q := quizCache.Quiz
	question.SetTitles(q.Title, briefSection, subSection)

	question.QuizUsesMathML = q.UsesMathML

	// Update the section and subSection,
	// so we can return it in the JSON,
	// so the caller of the REST API does not need to discover these details separately.
	if section != nil {
		question.SectionId = section.Id
		question.Section = &(section.HasIdAndTitle)
	}

	if subSection != nil {
		question.SubSectionId = subSection.Id
		question.SubSection = &(subSection.HasIdAndTitle)
	}

	return nil
}

func handleErrorAsHttpError(w http.ResponseWriter, code int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	log.Print(msg)

	http.Error(w, msg, code)
}

// marshalAndWriteOrHttpError() writes the object to the writer as JSON.
func marshalAndWriteOrHttpError(w http.ResponseWriter, v interface{}) {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "json.Marshal() failed: %v", err)
		return
	}

	_, err = w.Write(jsonStr)
	if err != nil {
		handleErrorAsHttpError(w, http.StatusInternalServerError, "w.Write() failed: %v", err)
	}
}
