package data

// ID is used as a key, so it must be comparable, type alias is used to allow plugging in different types
type ID string

type User struct {
	ID      ID
	Name    string
	Surname string
}
