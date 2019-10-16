package user

// Statistics for the user, either per quiz (SectionId is then empty), or per section in a quiz.
type Stats struct {
	QuizId    string
	SectionId string

	Answered int
	Correct  int

	CountQuestionsAnsweredOnce int
	CountQuestionsCorrectOnce  int

	QuestionHistories []QuestionHistory
}

func (self *Stats) GetQuestionCountAnsweredWrong(questionId string) int {
	qh, ok := self.getQuestionHistoryForQuestionId(questionId)
	if !ok {
		return 0
	}

	return qh.CountAnsweredWrong
}

func (self *Stats) GetQuestionWasAnswered(questionId string) bool {
	_, found := self.getQuestionHistoryForQuestionId(questionId)
	return found
}

func (self *Stats) IncrementAnswered() {
	self.Answered += 1
}

func (self *Stats) IncrementCorrect() {
	self.Correct += 1
}

func (self *Stats) getQuestionHistoryForQuestionId(questionId string) (*QuestionHistory, bool) {
	if self.QuestionHistories == nil {
		return nil, false
	}

	// TODO: Performance.
	// We would ideally use a map here,
	// but Go's datastore library does not allow that as an entity field type.
	for i := range self.QuestionHistories {
		qh := &self.QuestionHistories[i]
		if qh.QuestionId == questionId {
			return qh, true
		}
	}

	return nil, false
}

func (self *Stats) UpdateProblemQuestion(questionId string, answerIsCorrect bool) {
	if len(questionId) == 0 {
		// Log.error("updateProblemQuestion(): questionId is empty.");
		return
	}

	firstTimeAsked := false
	firstTimeCorrect := false

	questionHistory, exists := self.getQuestionHistoryForQuestionId(questionId)

	//Add a new one, if necessary:
	if !exists {
		firstTimeAsked = true
		if answerIsCorrect {
			firstTimeCorrect = true
		}

		questionHistory = new(QuestionHistory)
		questionHistory.QuestionId = questionId
	} else if answerIsCorrect && !questionHistory.AnsweredCorrectlyOnce {
		firstTimeCorrect = true
	}

	//Increase the wrong-answer count:
	questionHistory.AdjustCount(answerIsCorrect)

	if firstTimeAsked {
		self.CountQuestionsAnsweredOnce++
	}

	if firstTimeCorrect {
		self.CountQuestionsCorrectOnce++
	}

	if !exists {
		self.QuestionHistories = append(self.QuestionHistories, *questionHistory)
	}
	//TODO? cacheIsInvalid = true;
}

func (self *QuestionHistory) AdjustCount(result bool) {
	if result {
		self.AnsweredCorrectlyOnce = true
	}

	if result {
		self.CountAnsweredWrong -= 1
	} else {
		self.CountAnsweredWrong += 1
	}
}
