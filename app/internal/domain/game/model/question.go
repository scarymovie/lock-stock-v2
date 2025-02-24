package model

type Question struct {
	text  string
	hints []*Hint
}

func NewQuestion(text string, hints []*Hint) *Question {
	return &Question{text: text, hints: hints}
}
