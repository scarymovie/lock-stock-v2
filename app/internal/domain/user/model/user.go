package model

type User struct {
	uid  string
	name string
}

func NewUser(uid string, name string) *User {
	return &User{uid: uid, name: name}
}

func (u *User) Uid() string {
	return u.uid
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetUid(uid string) {
	u.uid = uid
}
func (u *User) SetName(name string) {
	u.name = name
}
