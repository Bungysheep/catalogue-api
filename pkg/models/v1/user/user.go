package user

import "time"

// User type
type User struct {
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Status     string    `json:"status"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedBy string    `json:"modified_by"`
	ModifiedAt time.Time `json:"modified_at"`
	Vers       int64     `json:"vers"`
}

// NewUser - Creates user
func NewUser() *User {
	return &User{}
}

// GetUsername - Returns username
func (usr *User) GetUsername() string {
	return usr.Username
}

// GetName - Returns name
func (usr *User) GetName() string {
	return usr.Name
}

// GetEmail - Returns email
func (usr *User) GetEmail() string {
	return usr.Email
}

// GetPassword - Returns password
func (usr *User) GetPassword() string {
	return usr.Password
}

// GetStatus - Returns status
func (usr *User) GetStatus() string {
	return usr.Status
}

// GetCreatedBy - Returns created by
func (usr *User) GetCreatedBy() string {
	return usr.CreatedBy
}

// GetCreatedAt - Returns created at
func (usr *User) GetCreatedAt() time.Time {
	return usr.CreatedAt
}

// GetModifiedBy - Returns modified by
func (usr *User) GetModifiedBy() string {
	return usr.ModifiedBy
}

// GetModifiedAt - Returns modified at
func (usr *User) GetModifiedAt() time.Time {
	return usr.ModifiedAt
}

// GetVers - Returns vers
func (usr *User) GetVers() int64 {
	return usr.Vers
}
