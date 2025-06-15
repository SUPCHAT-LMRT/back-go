package models

// PollResponse représente la réponse pour un sondage
type PollResponse struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Options   []Option   `json:"options"`
	CreatedAt string     `json:"created_at"`
	CreatedBy UserDetail `json:"created_by"`
}

// Option représente une option de sondage
type Option struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Votes int    `json:"votes"`
}

// UserDetail représente les informations d'un utilisateur
type UserDetail struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type KickMemberRequest struct {
	MemberId string `json:"member_id" example:"user-123"`
}

type RecentChatResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	LastMessage string       `json:"last_message"`
	UpdatedAt   string       `json:"updated_at"`
	Members     []UserDetail `json:"members"`
}

type AssignJobRequest struct {
	JobId    string `json:"job_id" example:"job-123"`
	UserId   string `json:"user_id" example:"user-456"`
	DueDate  string `json:"due_date,omitempty" example:"2025-07-15T14:00:00Z"`
	Priority int    `json:"priority,omitempty" example:"2"`
}

type CreateJobRequest struct {
	Name string `json:"name" example:"Développeur Backend"`
}

type CreateJobResponse struct {
	Id   string `json:"id" example:"job-123"`
	Name string `json:"name" example:"Développeur Backend"`
}

type ListJobsResponse struct {
	Jobs []JobResponse `json:"jobs"`
}

type JobResponse struct {
	Id   string `json:"id" example:"job-123"`
	Name string `json:"name" example:"Développeur Frontend"`
}

type CheckPermissionsRequest struct {
	UserId     string `json:"user_id" example:"user-123"`
	Permission string `json:"permission" example:"create_poll"`
}

type CheckPermissionsResponse struct {
	HasPermission bool   `json:"has_permission" example:"true"`
	Message       string `json:"message,omitempty" example:"L'utilisateur a la permission"`
}

type UnassignJobRequest struct {
	UserId string `json:"user_id" example:"user-123"`
	JobId  string `json:"job_id" example:"job-456"`
}

type UnassignJobResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"Rôle retiré avec succès"`
}

type UpdateJobRequest struct {
	Id   string `json:"id" example:"job-123" binding:"required"`
	Name string `json:"name" example:"Développeur Senior" binding:"required"`
}

type UpdateJobResponse struct {
	Id      string `json:"id" example:"job-123"`
	Name    string `json:"name" example:"Développeur Senior"`
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"Rôle mis à jour avec succès"`
}

type SaveStatusRequest struct {
	Status string `json:"status" example:"available" binding:"required" enums:"available,away,busy,offline"`
	Text   string `json:"text,omitempty" example:"En réunion jusqu'à 15h"`
}

// SaveStatusResponse représente la réponse après mise à jour du statut
type SaveStatusResponse struct {
	Success bool   `json:"success" example:"true"`
	Status  string `json:"status" example:"available"`
	Text    string `json:"text,omitempty" example:"En réunion jusqu'à 15h"`
}

type RequestForgotPasswordRequest struct {
	Email string `json:"email" example:"utilisateur@exemple.com" binding:"required,email"`
}

// RequestForgotPasswordResponse est la réponse à une demande de mot de passe oublié
type RequestForgotPasswordResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"Un email de réinitialisation a été envoyé"`
}

type WorkspaceMemberResponse struct {
	ID        string `json:"id" example:"user-123"`
	Email     string `json:"email" example:"utilisateur@exemple.com"`
	FirstName string `json:"first_name" example:"Jean"`
	LastName  string `json:"last_name" example:"Dupont"`
	Role      string `json:"role" example:"admin"`
	Status    string `json:"status" example:"active"`
	AvatarURL string `json:"avatar_url,omitempty" example:"https://exemple.com/avatar.jpg"`
}

// ListWorkspaceMembersResponse représente la liste des membres d'un espace de travail
type ListWorkspaceMembersResponse struct {
	Members []WorkspaceMemberResponse `json:"members"`
	Total   int                       `json:"total" example:"10"`
}
type ErrorResponse struct {
	Error string `json:"error" example:"Message d'erreur"`
}
