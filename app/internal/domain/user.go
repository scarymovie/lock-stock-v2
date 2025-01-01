package domain

type User struct {
	Id string
}

func (u User) GetUserId() string {
	return u.Id
}
