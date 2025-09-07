package constant

// Permission adalah tipe untuk permission user
type Permission string

const (
	PermissionRead   Permission = "read"
	PermissionWrite  Permission = "write"
	PermissionDelete Permission = "delete"
)

// String mengembalikan string dari Permission
func (p Permission) String() string {
	return string(p)
}
