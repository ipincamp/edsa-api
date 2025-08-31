package constant

type Permission string

const (
	PermissionRead   Permission = "read"
	PermissionWrite  Permission = "write"
	PermissionDelete Permission = "delete"
)

func (r Permission) String() string {
	return string(r)
}
