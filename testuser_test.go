package maildoor_test

type testUser string

func (tu testUser) EmailAddress() string {
	return string(tu)
}
