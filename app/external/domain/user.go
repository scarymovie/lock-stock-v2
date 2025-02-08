package domain

type User interface {
	Uid() string
	Name() string
	SetUid(uid string)
	SetName(name string)
}
