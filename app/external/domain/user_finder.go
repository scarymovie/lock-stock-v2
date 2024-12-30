package domain

type UserFinder interface {
	FindById(userId string) (User, error)
}
