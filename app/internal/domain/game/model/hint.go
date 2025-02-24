package model

type Hint struct {
	text string
}

func NewHint(text string) *Hint {
	return &Hint{text: text}
}
