package constant

// Role adalah tipe untuk role user
type Role string

const (
	RoleAdmin   Role = "admin"
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
	RoleGuest   Role = "guest"
)

// String mengembalikan string dari Role
func (r Role) String() string {
	return string(r)
}
