package configs

const (
	// SHORTDATEFORMAT - Short Date format
	SHORTDATEFORMAT = "2006-01-02"

	// DATEFORMAT - Date format
	DATEFORMAT = "2006-01-02 15:04:05"

	// DEFAULTDATE - Default Date
	DEFAULTDATE = "2000-01-01 00:00:00"

	// TESTDOMAIN - Test domain
	TESTDOMAIN = "http://localhost:50051"

	// PORT - Port number
	PORT = "50051"

	// READTIMEOUT - Read timeout
	READTIMEOUT = 10

	// WRITETIMEOUT - Write timeout
	WRITETIMEOUT = 10

	// PGTESTCONNSTRING - Postgres connection string
	PGTESTCONNSTRING = "postgres://clg_local_dev:clg_local_dev@localhost:5432/clg_local_dev?sslmode=disable"

	// TOKENLIFETIME - Jwt token lifetime
	TOKENLIFETIME = 5

	// TOKENSIGNKEY - Jwt token sign in key
	TOKENSIGNKEY = "secret"
)
