package seeders

import (
	"fmt"
	"log"
	"strings"

	"github.com/ipincamp/go-edsa-api/internal/database/factories"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"gorm.io/gorm"
)

func ClassGroupSeeder(db *gorm.DB, logger *log.Logger) error {
	logger.Println("Seeding subjects, classes, groups, and user_groups...")

	// 1. Dapatkan Role ID
	var teacherRole repo.RoleGORM
	if err := db.Where("name = ?", domain.RoleNameTeacher).First(&teacherRole).Error; err != nil {
		return fmt.Errorf("failed to find 'teacher' role: %w", err)
	}
	var studentRole repo.RoleGORM
	if err := db.Where("name = ?", domain.RoleNameStudent).First(&studentRole).Error; err != nil {
		return fmt.Errorf("failed to find 'student' role: %w", err)
	}

	// 2. Hash password default SEKALI di sini
	passSvc := argon2id.NewPasswordService()
	defaultPassword := "password"
	hashedDefaultPassword, err := passSvc.Hash(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed to hash default password for ClassGroupSeeder: %w", err)
	}
	logger.Println("   Hashed default password for teachers and students in groups.")

	// 3. Ambil email yang sudah ada
	existingEmails, err := getExistingEmails(db)
	if err != nil {
		return err
	}
	logger.Printf("   Loaded %d existing emails for ClassGroupSeeder.", len(existingEmails))

	// 4. Data subject
	subjectData := map[string]int{
		"English for Beginner":     5,
		"English for Intermediate": 3,
		"English for Expert":       1,
	}

	// 5. Iterasi
	for subjectName, numGroups := range subjectData {
		// --- BUAT SUBJECT ---
		subject := repo.SubjectGORM{Name: subjectName}
		db.FirstOrCreate(&subject, repo.SubjectGORM{Name: subject.Name})
		logger.Printf("Seeded Subject: %s (ID: %d)", subject.Name, subject.ID)

		// --- BUAT CLASS ---
		// Cth: "Beginner Class", "Intermediate Class"
		className := strings.Split(subject.Name, " ")[1] + " Class"
		class := repo.ClassGORM{Name: className, SubjectID: subject.ID}
		db.FirstOrCreate(&class, repo.ClassGORM{Name: class.Name, SubjectID: class.SubjectID})
		logger.Printf("  -> Seeded Class: %s (ID: %d)", class.Name, class.ID)

		// --- BUAT GROUP ---
		for i := 1; i <= numGroups; i++ {
			// Cth: "Beginner Group 1", "Beginner Group 2", ...
			groupName := fmt.Sprintf("%s Group %d", strings.Split(subject.Name, " ")[1], i)
			group := repo.GroupGORM{Name: groupName, ClassID: class.ID}
			db.FirstOrCreate(&group, repo.GroupGORM{Name: group.Name, ClassID: class.ID})
			logger.Printf("    -> Seeded Group: %s (ID: %d)", group.Name, group.ID)

			// --- BUAT USERS (TEACHER & STUDENT) UNTUK GROUP INI ---
			const numTeachers = 2
			const numStudents = 40
			usersToAssign := make([]*repo.UserGORM, 0, numTeachers+numStudents)
			generatedEmails := make(map[string]bool)
			attempts := 0
			maxAttempts := (numTeachers + numStudents) * 5

			// Buat Teacher
			for len(usersToAssign) < numTeachers && attempts < maxAttempts {
				attempts++
				// Panggil factory dengan hash default
				teacher := factories.UserFactory(db, teacherRole.ID, "avatar3.png", hashedDefaultPassword)
				if existingEmails[teacher.Email] || generatedEmails[teacher.Email] {
					continue
				}
				usersToAssign = append(usersToAssign, teacher)
				generatedEmails[teacher.Email] = true
				existingEmails[teacher.Email] = true
			}

			// Buat Student
			targetUserCount := numTeachers + numStudents
			for len(usersToAssign) < targetUserCount && attempts < maxAttempts {
				attempts++
				// Panggil factory dengan hash default
				student := factories.UserFactory(db, studentRole.ID, "avatar2.png", hashedDefaultPassword)
				if existingEmails[student.Email] || generatedEmails[student.Email] {
					continue
				}
				usersToAssign = append(usersToAssign, student)
				generatedEmails[student.Email] = true
				existingEmails[student.Email] = true
			}

			if len(usersToAssign) < targetUserCount {
				logger.Printf("      WARNING: Could only generate %d unique users for group %s after %d attempts.", len(usersToAssign), group.Name, attempts)
			}

			// Batch Insert Users & Assign to Group
			if len(usersToAssign) > 0 {
				logger.Printf("      -> Inserting %d users for group %s...", len(usersToAssign), group.Name)
				if err := db.Create(&usersToAssign).Error; err != nil {
					return fmt.Errorf("failed to batch insert users for group %s: %w", group.Name, err)
				}
				logger.Printf("      -> Successfully inserted %d users.", len(usersToAssign))

				if err := db.Model(&group).Association("Users").Append(usersToAssign); err != nil {
					return fmt.Errorf("failed to assign users to group %s: %w", group.Name, err)
				}
				logger.Printf("      -> Assigned %d users to group %s", len(usersToAssign), group.Name)
			} else {
				logger.Printf("      -> No new unique users generated for group %s", group.Name)
			}
		}
	}

	logger.Println("Finished seeding classes and groups.")
	return nil
}
