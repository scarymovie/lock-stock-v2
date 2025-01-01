package domain

type User struct {
	Id  int
	Uid string
}

func (u User) GetUserId() int {
	return u.Id
}

func (u User) GetUserUid() string {
	return u.Uid
}
