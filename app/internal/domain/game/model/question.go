package model

type Question struct {
	text  string
	hints []*Hint
}

func (q *Question) Text() string {
	return q.text
}

func (q *Question) SetText(text string) {
	q.text = text
}

func (q *Question) Hints() []*Hint {
	return q.hints
}

func (q *Question) SetHints(hints []*Hint) {
	q.hints = hints
}

func NewQuestion(text string, hints []*Hint) *Question {
	return &Question{text: text, hints: hints}
}
