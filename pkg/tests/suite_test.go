package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/configs"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/signinclaimresource"
	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
	"github.com/bungysheep/catalogue-api/pkg/protocols/rest"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

var accessTokenTest string

func TestMain(m *testing.M) {
	ctx := context.TODO()

	restServer := &rest.Server{}

	if err := database.CreateDbConnection(); err != nil {
		log.Printf("Failed create database connection, error: %v.\n", err)
		os.Exit(1)
	}

	go func() {
		restServer.RunServer()
	}()

	defer func() {
		ctx, err := context.WithTimeout(context.TODO(), 5*time.Second)
		if err != nil {
			log.Printf("Failed create context timeout, error: %v.\n", err)
			os.Exit(1)
		}
		restServer.Server.Shutdown(ctx)
	}()

	tearUp(ctx)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func tearUp(ctx context.Context) {
	setupDatabase(ctx)

	setupAuthUser()
}

func setupDatabase(ctx context.Context) {
	tx, _ := database.DbConnection.Begin()

	var err error
	// Truncate all tables
	_, err = tx.Exec("TRUNCATE TABLE users")
	_, err = tx.Exec("TRUNCATE TABLE catalogues")

	// Seed users
	_, err = tx.Exec(`INSERT INTO users (username, name, email, password, status, created_by, created_at, modified_by, modified_at, vers) VALUES 
					('TESTUSER', 'Test User', 'Test.User@testmail.com', '$2a$10$6zFv6A/AzTEUxbmKGnBOAOhwksvcnopCQelGMskeyT1z8ONFswLzy', 'I', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1)`)

	// Seed catalogues
	_, err = tx.Exec(`INSERT INTO catalogues (code, descr, details, created_by, created_at, modified_by, modified_at, vers) VALUES 
					('CLG_TEST_1', 'Catalogue Test 1', 'Catalogue Test 1', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_2', 'Catalogue Test 2', 'Catalogue Test 2', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_3', 'Catalogue Test 3', 'Catalogue Test 3', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1)`)

	if err != nil {
		tx.Rollback()
	}

	tx.Commit()
}

func setupAuthUser() {
	expiresAt := time.Now().Add(15 * time.Minute).Unix()
	signInToken := &signinclaimresource.SignInClaimResource{
		Username: "TESTUSER",
		Name:     "Test User",
		Email:    "Test.User@testmail.com",
		Status:   "A",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, signInToken)
	signedToken, _ := token.SignedString([]byte(configs.TOKENSIGNKEY))
	accessTokenTest = "Bearer " + signedToken
}