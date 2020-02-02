package authcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/commons/status"
	"github.com/bungysheep/catalogue-api/pkg/configs"
	"github.com/bungysheep/catalogue-api/pkg/controllers/v1/basecontroller"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/signinclaimresource"
	usermodel "github.com/bungysheep/catalogue-api/pkg/models/v1/user"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/userrepository"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// AuthController type
type AuthController struct {
	basecontroller.BaseResource
}

// SignInResponseResource type
type SignInResponseResource struct {
	Username string                 `json:"username"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	Status   string                 `json:"status"`
	Token    map[string]interface{} `json:"token"`
}

// NewAuthController - Creates auth controller
func NewAuthController() *AuthController {
	return &AuthController{}
}

// GetAll - Return all catalogues
func (authCtl *AuthController) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("Retrieving all users.\n")

	usrRepo := userrepository.NewUserRepository()
	result, err := usrRepo.GetAll(r.Context())
	if err != nil {
		authCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	authCtl.WriteResponse(w, http.StatusOK, true, result, "")
}

// SignIn - Sign in user and return access token
func (authCtl *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	log.Printf("Sign in user.\n")

	user := usermodel.NewUser()
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		authCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Invalid sign in request.")
		return
	}

	usrRepo := userrepository.NewUserRepository()
	result, err := usrRepo.GetByUsername(r.Context(), user.GetUsername())
	if err != nil {
		authCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if result == nil || result.GetStatus() != "A" {
		authCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Incorrect username/password.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.GetPassword()), []byte(user.Password)); err != nil {
		authCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Incorrect username/password.")
		return
	}

	expiresAt := time.Now().Add(configs.TOKENLIFETIME * time.Minute)
	tokenClaims := signinclaimresource.NewSignInClaimResource()
	tokenClaims.Username = result.GetUsername()
	tokenClaims.Name = result.GetName()
	tokenClaims.Email = result.GetEmail()
	tokenClaims.Status = result.GetStatus()
	tokenClaims.StandardClaims = &jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString([]byte(configs.TOKENSIGNKEY))
	if err != nil {
		authCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Failed to sign token.")
		return
	}

	signInResp := &SignInResponseResource{
		Username: result.GetUsername(),
		Name:     result.GetName(),
		Email:    result.GetEmail(),
		Status:   result.GetStatus(),
		Token: map[string]interface{}{
			"access_token": tokenString,
			"expires_at":   expiresAt.Format(time.RFC3339),
		},
	}

	authCtl.WriteResponse(w, http.StatusOK, true, signInResp, "")
}

// Register - Register user
func (authCtl *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Registering user.\n")

	newUsr := usermodel.NewUser()
	err := json.NewDecoder(r.Body).Decode(newUsr)
	if err != nil {
		authCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Invalid register request.")
		return
	}

	valid, message := newUsr.DoValidate()
	if !valid {
		authCtl.WriteResponse(w, http.StatusBadRequest, false, nil, message)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(newUsr.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		authCtl.WriteResponse(w, http.StatusBadRequest, false, nil, "Failed to encrypt password.")
		return
	}

	newUsr.Password = string(pass)
	newUsr.Status = status.Active.String()
	newUsr.CreatedBy = newUsr.GetUsername()
	newUsr.ModifiedBy = newUsr.GetUsername()
	newUsr.Vers = 1

	usrRepo := userrepository.NewUserRepository()
	nbrRows, err := usrRepo.Create(r.Context(), newUsr)
	if err != nil {
		authCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	if nbrRows == 0 {
		authCtl.WriteResponse(w, http.StatusNotFound, false, nil, "User was not registered.")
		return
	}

	result, err := usrRepo.GetByUsername(r.Context(), newUsr.GetUsername())
	if err != nil {
		authCtl.WriteResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	authCtl.WriteResponse(w, http.StatusAccepted, true, result, "User has been registered.")
}
