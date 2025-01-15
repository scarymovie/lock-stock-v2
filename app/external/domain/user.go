package domain

type User interface {
	GetUserId() int
	GetUserUid() string
	GetUserName() string
	SetUserUid(uid string)
	SetUserName(name string)
}
