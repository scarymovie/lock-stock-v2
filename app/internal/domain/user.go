package domain

type User struct {
	Id string
}

func (u User) GetId() string {
	return u.Id
}
