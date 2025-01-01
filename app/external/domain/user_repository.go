package domain

type UserRepository interface {
	SaveUser(user User) error
}
