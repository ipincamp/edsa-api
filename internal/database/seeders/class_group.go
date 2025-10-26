package seeders

import (
	"fmt"
	"log"
	"strings"

	"github.com/ipincamp/go-edsa-api/internal/database/factories"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

func ClassGroupSeeder(db *gorm.DB) error {
	log.Println("Seeding subjects, classes, groups, and user_groups...")

	// 1. Dapatkan ID role "teacher" dan "student"
	var teacherRole repo.RoleGORM
	if err := db.Where("name = ?", domain.RoleNameTeacher).First(&teacherRole).Error; err != nil {
		return fmt.Errorf("failed to find 'teacher' role: %w", err)
	}
	var studentRole repo.RoleGORM
	if err := db.Where("name = ?", domain.RoleNameStudent).First(&studentRole).Error; err != nil {
		return fmt.Errorf("failed to find 'student' role: %w", err)
	}

	// 2. Ambil semua email yang sudah ada SEBELUM loop utama
	existingEmails, err := getExistingEmails(db)
	if err != nil {
		return err
	}
	log.Printf("Loaded %d existing emails for ClassGroupSeeder.", len(existingEmails))

	// 3. Tentukan data subject
	subjectData := map[string]int{
		"English for Beginner":     5,
		"English for Intermediate": 3,
		"English for Expert":       1,
	}

	// 3. Mulai iterasi
	for subjectName, numGroups := range subjectData {
		// --- BUAT SUBJECT ---
		subject := repo.SubjectGORM{Name: subjectName}
		if err := db.FirstOrCreate(&subject, repo.SubjectGORM{Name: subject.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed subject '%s': %w", subject.Name, err)
		}
		if subject.ID != 0 {
			log.Printf("Seeded Subject: %s (ID: %d)", subject.Name, subject.ID)
		}

		// --- BUAT 1 CLASS PER SUBJECT ---
		// Cth: "Beginner Class", "Intermediate Class"
		className := strings.Split(subject.Name, " ")[1] + " Class"
		class := repo.ClassGORM{Name: className, SubjectID: subject.ID}
		if err := db.FirstOrCreate(&class, repo.ClassGORM{Name: class.Name, SubjectID: class.SubjectID}).Error; err != nil {
			return fmt.Errorf("failed to seed class '%s': %w", class.Name, err)
		}
		if class.ID != 0 {
			log.Printf("  -> Seeded Class: %s (ID: %d)", class.Name, class.ID)
		}

		// --- BUAT N GROUPS PER CLASS ---
		for i := 1; i <= numGroups; i++ {
			// Cth: "Beginner Group 1", "Beginner Group 2", ...
			groupName := fmt.Sprintf("%s Group %d", strings.Split(subject.Name, " ")[1], i)
			group := repo.GroupGORM{Name: groupName, ClassID: class.ID}
			db.FirstOrCreate(&group, repo.GroupGORM{Name: group.Name, ClassID: class.ID})
			log.Printf("    -> Seeded Group: %s (ID: %d)", group.Name, group.ID)

			// --- BUAT USERS (TEACHER & STUDENT) UNTUK GROUP INI ---
			const numTeachers = 2
			const numStudents = 40
			usersToAssign := make([]*repo.UserGORM, 0, numTeachers+numStudents)
			generatedEmails := make(map[string]bool) // Email unik dalam batch group ini
			attempts := 0
			maxAttempts := (numTeachers + numStudents) * 5 // Toleransi percobaan

			// Buat Teacher
			for len(usersToAssign) < numTeachers && attempts < maxAttempts {
				attempts++
				teacher := factories.UserFactory(db, teacherRole.ID, "avatar3.png")
				// Cek duplikasi (global & batch)
				if existingEmails[teacher.Email] || generatedEmails[teacher.Email] {
					continue
				}
				usersToAssign = append(usersToAssign, teacher)
				generatedEmails[teacher.Email] = true
				existingEmails[teacher.Email] = true // Tambahkan ke global set agar tidak dipakai group lain
			}

			// Buat Student
			targetUserCount := numTeachers + numStudents
			for len(usersToAssign) < targetUserCount && attempts < maxAttempts {
				attempts++
				student := factories.UserFactory(db, studentRole.ID, "avatar2.png")
				// Cek duplikasi (global & batch)
				if existingEmails[student.Email] || generatedEmails[student.Email] {
					continue
				}
				usersToAssign = append(usersToAssign, student)
				generatedEmails[student.Email] = true
				existingEmails[student.Email] = true // Tambahkan ke global set
			}

			// Laporkan jika gagal membuat semua user
			if len(usersToAssign) < targetUserCount {
				log.Printf("      WARNING: Could only generate %d unique users for group %s after %d attempts.", len(usersToAssign), group.Name, attempts)
			}

			// Batch Insert Users ke DB
			if len(usersToAssign) > 0 {
				log.Printf("      -> Inserting %d users for group %s...", len(usersToAssign), group.Name)
				if err := db.Create(&usersToAssign).Error; err != nil {
					// Jika batch insert gagal, log error dan mungkin hentikan seeder
					return fmt.Errorf("failed to batch insert users for group %s: %w", group.Name, err)
				}
				log.Printf("      -> Successfully inserted %d users.", len(usersToAssign))

				// Tambahkan relasi User-Group (Append association)
				if err := db.Model(&group).Association("Users").Append(usersToAssign); err != nil {
					return fmt.Errorf("failed to assign users to group %s: %w", group.Name, err)
				}
				log.Printf("      -> Assigned %d users to group %s", len(usersToAssign), group.Name)
			} else {
				log.Printf("      -> No new unique users generated for group %s", group.Name)
			}
		}
	}

	log.Println("Finished seeding classes and groups.")
	return nil
}
