package sample

type User string

func (u User) EmailAddress() string {
	return string(u)
}
