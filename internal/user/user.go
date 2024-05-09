package user

type User struct {
	ID             string
	UserType       UserType
	PhoneNumber    string
	Name           string
	HashedPassword string
}
type UserType string

const (
	Staff    UserType = "Staff"
	Customer UserType = "Customer"
)
