package constant

// Permission adalah tipe custom untuk semua hak akses di dalam sistem.
type Permission string

// String mengembalikan representasi string dari Permission.
func (p Permission) String() string {
	return string(p)
}

// --- Daftar Konstanta Permission ---
const (
	// User Permissions
	UsersViewSelf    Permission = "users.view.self"    // Melihat profil diri sendiri
	UsersViewOther   Permission = "users.view.other"   // Melihat profil pengguna lain
	UsersListAll     Permission = "users.list.all"     // Melihat daftar semua pengguna (paginasi)
	UsersUpdateSelf  Permission = "users.update.self"  // Memperbarui profil diri sendiri
	UsersUpdateOther Permission = "users.update.other" // Memperbarui profil pengguna lain
	UsersDelete      Permission = "users.delete"       // Menghapus pengguna
	UsersCreate      Permission = "users.create"       // Membuat pengguna baru

	// Role & Permission Management
	RolesManage Permission = "roles.manage" // Mengelola role dan permission

	// Course Permissions (Contoh)
	CoursesCreate Permission = "courses.create"
	CoursesUpdate Permission = "courses.update"
	CoursesRead   Permission = "courses.read"

	// Grade Permissions (Contoh)
	GradesManage Permission = "grades.manage"

	// Assignment Permissions (Contoh)
	AssignmentSubmit Permission = "assignment.submit"
)
