package data

type Repository interface {
	Get(ID) (*User, error)
}
