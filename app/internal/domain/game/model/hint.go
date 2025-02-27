package model

type Hint struct {
	text string
}

func (h *Hint) Text() string {
	return h.text
}

func (h *Hint) SetText(text string) {
	h.text = text
}

func NewHint(text string) *Hint {
	return &Hint{text: text}
}
