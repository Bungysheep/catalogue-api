package user

import (
	"github.com/bungysheep/catalogue-api/pkg/models/v1/basemodel"
	"strings"
	"time"
)

// User type
type User struct {
	basemodel.BaseModel
	Username   string    `json:"username" mandatory:"true" max_length:"16"`
	Name       string    `json:"name" mandatory:"true" max_length:"64"`
	Email      string    `json:"email" mandatory:"true" max_length:"255"`
	Password   string    `json:"password" mandatory:"true" max_length:"16"`
	Status     string    `json:"status" mandatory:"true" max_length:"1"`
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
	return strings.ToUpper(usr.Username)
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
	return strings.ToUpper(usr.Status)
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

// DoValidate - Validate user
func (usr *User) DoValidate() (bool, string) {
	return usr.DoValidateBase(*usr)
}
