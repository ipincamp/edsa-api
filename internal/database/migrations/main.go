package migrations

import "github.com/go-gormigrate/gormigrate/v2"

// Kumpulkan semua migrasi di sini
func GetAllMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		CreateRolesTable(),
		CreateUsersTable(),
		CreateSubjectsTable(),
		CreateClassesTable(),
		CreateGroupsTable(),
		CreateUserGroupsTable(),
		CreateBooksTable(),
		CreatePagesTable(),
		CreateInteractionsTable(),
		CreateUserBookProgressTable(),
		CreateGamesTable(),
		CreateUserGameScoresTable(),
		CreateActivityLogsTable(),
		CreateMediaAssetsTable(),
		// AnotherMigration(),
	}
}
