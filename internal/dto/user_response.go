package dto

type UserResponse struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	FullName    string   `json:"full_name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions,omitempty"`
}
