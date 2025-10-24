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

	// 2. Tentukan data yang akan di-seed
	// Format: map[NamaSubject]jumlahGrup
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
			if err := db.FirstOrCreate(&group, repo.GroupGORM{Name: group.Name, ClassID: class.ID}).Error; err != nil {
				return fmt.Errorf("failed to seed group '%s': %w", group.Name, err)
			}
			if group.ID != 0 {
				log.Printf("    -> Seeded Group: %s (ID: %d)", group.Name, group.ID)
			}

			// --- BUAT USERS DAN MASUKKAN KE GROUP (USER_GROUPS) ---
			var usersToAssign []*repo.UserGORM

			// Buat 2 Teacher baru untuk grup ini
			for j := 0; j < 2; j++ {
				teacher := factories.UserFactory(teacherRole.ID)

				// Cek email unik (penting!)
				var existing repo.UserGORM
				if db.Where("email = ?", teacher.Email).First(&existing).Error == nil {
					j-- // Coba lagi jika email sudah ada
					continue
				}

				if err := db.Create(teacher).Error; err != nil {
					return fmt.Errorf("failed to create teacher for group %s: %w", group.Name, err)
				}
				usersToAssign = append(usersToAssign, teacher)
			}

			// Buat 40 Student baru untuk grup ini
			for j := 0; j < 40; j++ {
				student := factories.UserFactory(studentRole.ID)

				// Cek email unik
				var existing repo.UserGORM
				if db.Where("email = ?", student.Email).First(&existing).Error == nil {
					j-- // Coba lagi jika email sudah ada
					continue
				}

				if err := db.Create(student).Error; err != nil {
					return fmt.Errorf("failed to create student for group %s: %w", group.Name, err)
				}
				usersToAssign = append(usersToAssign, student)
			}

			// Masukkan semua user baru (2T + 40S) ke tabel relasi 'user_groups'
			if err := db.Model(&group).Association("Users").Append(usersToAssign); err != nil {
				return fmt.Errorf("failed to assign users to group %s: %w", group.Name, err)
			}
			log.Printf("      -> Created and assigned 2 teachers and 40 students to group %s", group.Name)
		}
	}

	log.Println("Finished seeding classes and groups.")
	return nil
}
