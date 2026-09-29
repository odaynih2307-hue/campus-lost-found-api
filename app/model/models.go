package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
type Item struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	ReportedBy  int64     `json:"reported_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Claim struct {
	ID         int64     `json:"id"`
	ItemID     int64     `json:"item_id"`
	ClaimantID int64     `json:"claimant_id"`
	Note       string    `json:"note"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type CreateItemRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=120"`
	Description string `json:"description" validate:"required,min=5"`
	Category    string `json:"category" validate:"required,min=2,max=60"`
	Location    string `json:"location" validate:"required,min=2,max=120"`
	Status      string `json:"status" validate:"required,oneof=lost found resolved"`
}
type UpdateItemRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=120"`
	Description string `json:"description" validate:"required,min=5"`
	Category    string `json:"category" validate:"required,min=2,max=60"`
	Location    string `json:"location" validate:"required,min=2,max=120"`
	Status      string `json:"status" validate:"required,oneof=lost found resolved"`
}
type PatchItemRequest struct {
	Title       *string `json:"title" validate:"omitnil,min=3,max=120"`
	Description *string `json:"description" validate:"omitnil,min=5"`
	Category    *string `json:"category" validate:"omitnil,min=2,max=60"`
	Location    *string `json:"location" validate:"omitnil,min=2,max=120"`
	Status      *string `json:"status" validate:"omitnil,oneof=lost found resolved"`
}
type CreateClaimRequest struct {
	Note string `json:"note" validate:"required,min=5"`
}
type UpdateClaimStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=approved rejected"`
}
type ListQuery struct {
	Limit                    int
	CursorCreated            *time.Time
	CursorID                 int64
	Search, Status, Category string
	OwnerID                  int64
}
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
type AuthUser struct {
	UserID         int64
	Username, Role string
}
