package constant

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
	RoleGuest   Role = "guest"
)

func (r Role) String() string {
	return string(r)
}
