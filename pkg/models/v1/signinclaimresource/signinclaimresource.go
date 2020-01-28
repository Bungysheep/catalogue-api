package signinclaimresource

import "github.com/dgrijalva/jwt-go"

// SignInClaimResource type
type SignInClaimResource struct {
	Username            string `json:"username"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	Status              string `json:"status"`
	*jwt.StandardClaims `json:"standard_claims"`
}

// NewSignInClaimResource - Creates sign in claim resource
func NewSignInClaimResource() *SignInClaimResource {
	return &SignInClaimResource{}
}

// GetUsername - Returns username
func (claim *SignInClaimResource) GetUsername() string {
	return claim.Username
}

// GetName - Returns name
func (claim *SignInClaimResource) GetName() string {
	return claim.Name
}

// GetEmail - Returns email
func (claim *SignInClaimResource) GetEmail() string {
	return claim.Email
}

// GetStatus - Returns status
func (claim *SignInClaimResource) GetStatus() string {
	return claim.Status
}
