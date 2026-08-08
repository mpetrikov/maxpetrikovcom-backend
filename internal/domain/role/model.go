package role

type Name string

const (
	Student Name = "student"
	Admin   Name = "admin"
)

type Role struct {
	ID   int16
	Name Name
}
