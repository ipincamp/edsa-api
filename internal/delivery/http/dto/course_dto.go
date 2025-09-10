package dto

import "github.com/ipincamp/go-edsa-api/internal/domain"

// ApplyToGroupRequest adalah DTO untuk request body saat guest mendaftar ke kelas.
type ApplyToGroupRequest struct {
	GroupCode string `json:"group_code" validate:"required"`
}

// HandleJoinRequest adalah DTO untuk request body saat guru memproses permintaan bergabung.
type HandleJoinRequest struct {
	// Digunakan pointer agar validator bisa membedakan antara nilai 'false' yang eksplisit
	// dan field yang tidak ada sama sekali. Namun untuk Fiber, `bool` biasa cukup.
	Approved bool `json:"approved" validate:"boolean"`
}

// CourseGroupResponse adalah DTO untuk menampilkan data kelas.
type CourseGroupResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	GroupCode        string `json:"group_code"`
	GroupDescription string `json:"group_description"`
	CourseName       string `json:"course_name,omitempty"` // Tampilkan nama mata pelajarannya
}

// ToCourseGroupResponse mengubah satu objek domain.CourseGroup menjadi dto.CourseGroupResponse.
func ToCourseGroupResponse(group domain.CourseGroup) CourseGroupResponse {
	response := CourseGroupResponse{
		ID:               group.ID,
		Name:             group.Name,
		GroupCode:        group.GroupCode,
		GroupDescription: group.GroupDescription,
	}
	// Cek jika relasi Course di-preload untuk menghindari panic error
	if group.Course.ID != "" {
		response.CourseName = group.Course.Name
	}
	return response
}

// ToCourseGroupListResponse mengubah slice dari domain.CourseGroup menjadi slice dto.CourseGroupResponse.
func ToCourseGroupListResponse(groups []domain.CourseGroup) []CourseGroupResponse {
	var response []CourseGroupResponse
	for _, group := range groups {
		response = append(response, ToCourseGroupResponse(group))
	}
	return response
}

type CourseRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateCourseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CourseResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func ToCourseResponse(course domain.Course) CourseResponse {
	return CourseResponse{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		CreatedAt:   course.CreatedAt.String(),
		UpdatedAt:   course.UpdatedAt.String(),
	}
}

func ToCourseListResponse(courses []domain.Course) []CourseResponse {
	var response []CourseResponse
	for _, course := range courses {
		response = append(response, ToCourseResponse(course))
	}
	return response
}
