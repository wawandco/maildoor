package sample

type User string

// EmailAddress returns the string for this sample user being used.
// This is the implementation of the needed Emailable interface method.
func (u User) EmailAddress() string {
	return string(u)
}
