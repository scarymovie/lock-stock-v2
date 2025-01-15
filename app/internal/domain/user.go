package domain

type User struct {
	Id   int
	Uid  string
	Name string
}

func (u User) GetUserId() int {
	return u.Id
}

func (u User) GetUserUid() string {
	return u.Uid
}

func (u User) GetUserName() string {
	return u.Name
}
