package dto

import "github.com/ipincamp/go-edsa-api/internal/domain"

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
