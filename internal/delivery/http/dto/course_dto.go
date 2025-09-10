package dto

import (
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
)

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

// CourseIDRequest adalah DTO untuk request yang memerlukan ID mata pelajaran di path parameter.
type CourseIDRequest struct {
	CourseId string `params:"courseId" validate:"required,uuid4"`
}

// CourseFilterRequest adalah DTO untuk filter list mata pelajaran.
type CourseFilterRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}

// CourseRequest adalah DTO untuk membuat mata pelajaran.
type CreateCourseRequest struct {
	Name        string `json:"name,omitempty" validate:"required,min=3,max=40"`
	Description string `json:"description,omitempty" validate:"omitempty,max=300"`
}

// UpdateCourseRequest adalah DTO untuk memperbarui mata pelajaran.
type UpdateCourseRequest = CreateCourseRequest

// CourseResponse adalah DTO untuk menampilkan detail mata pelajaran.
type CourseResponse struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// CourseListResponse adalah DTO untuk menampilkan daftar mata pelajaran.
type CourseListResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// ToCourseResponse mengubah satu objek domain.Course menjadi dto.CourseResponse.
func ToCourseResponse(course domain.Course) CourseResponse {
	return CourseResponse{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		CreatedAt:   course.CreatedAt.Format(constant.TimeFormat),
		UpdatedAt:   course.UpdatedAt.Format(constant.TimeFormat),
	}
}

// ToCourseListResponse mengubah slice dari domain.Course menjadi slice dto.CourseListResponse.
func ToCourseListResponse(courses []domain.Course) []CourseListResponse {
	var response []CourseListResponse
	for _, course := range courses {
		response = append(response, CourseListResponse{
			ID:        course.ID,
			Name:      course.Name,
			CreatedAt: course.CreatedAt.Format(constant.TimeFormat),
		})
	}
	return response
}
