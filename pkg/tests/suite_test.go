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
	_, err = tx.Exec("TRUNCATE TABLE custom_field_definitions")
	_, err = tx.Exec("TRUNCATE TABLE products")
	_, err = tx.Exec("TRUNCATE TABLE product_uoms")
	_, err = tx.Exec("TRUNCATE TABLE product_custom_fields")

	// Seed users
	_, err = tx.Exec(`INSERT INTO users (username, name, email, password, status, created_by, created_at, modified_by, modified_at, vers) VALUES 
					('TESTUSER', 'Test User', 'Test.User@testmail.com', '$2a$10$6zFv6A/AzTEUxbmKGnBOAOhwksvcnopCQelGMskeyT1z8ONFswLzy', 'I', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1)`)

	// Seed catalogues
	_, err = tx.Exec(`INSERT INTO catalogues (code, descr, details, created_by, created_at, modified_by, modified_at, vers) VALUES 
					('CLG_TEST_1', 'Catalogue Test 1', 'Catalogue Test 1', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_2', 'Catalogue Test 2', 'Catalogue Test 2', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_3', 'Catalogue Test 3', 'Catalogue Test 3', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1)`)

	// Seed custom field definitions
	_, err = tx.Exec(`INSERT INTO custom_field_definitions (clg_code, caption, type, mandatory, created_by, created_at, modified_by, modified_at, vers) VALUES 
					('CLG_TEST_1', 'Field-1', 'A', 1, 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_1', 'Field-2', 'N', 0, 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_1', 'Field-3', 'D', 0, 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_2', 'Field-1', 'A', 1, 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_2', 'Field-2', 'N', 0, 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_2', 'Field-3', 'D', 0, 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1)`)

	// Seed products
	_, err = tx.Exec(`INSERT INTO products (clg_code, code, descr, details, created_by, created_at, modified_by, modified_at, vers) VALUES 
					('CLG_TEST_1', 'P-0001', 'Book', 'Book', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_1', 'P-0002', 'Pen', 'Pen', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1),
					('CLG_TEST_1', 'P-0003', 'Tissue', 'Tissue', 'TESTUSER', CURRENT_TIMESTAMP, 'TESTUSER', CURRENT_TIMESTAMP, 1)`)

	// Seed product uoms
	_, err = tx.Exec(`INSERT INTO product_uoms (prod_id, code, descr, ratio, vers) VALUES 
					(1, 'EACH', 'Each', 1, 1),
					(1, 'BOX', 'Box', 2, 1),
					(2, 'EACH', 'Each', 1, 1),
					(2, 'BOX', 'Box', 2, 1),
					(3, 'EACH', 'Each', 1, 1),
					(3, 'PACK', 'Pack', 6, 1)`)

	// Seed product custom fields
	_, err = tx.Exec(`INSERT INTO product_custom_fields (prod_id, field_id, alpha_value, numeric_value, date_value) VALUES 
					(1, 1, 'Field-Prod-1', DEFAULT, DEFAULT),
					(1, 2, DEFAULT, 10.5, DEFAULT),
					(1, 3, DEFAULT, DEFAULT, '2020-01-01'),
					(2, 1, 'Field-Prod-2', DEFAULT, DEFAULT),
					(2, 2, DEFAULT, 20.5, DEFAULT),
					(2, 3, DEFAULT, DEFAULT, '2020-02-01'),
					(3, 1, 'Field-Prod-3', DEFAULT, DEFAULT),
					(3, 2, DEFAULT, 30.5, DEFAULT),
					(3, 3, DEFAULT, DEFAULT, '2020-03-01')`)

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
