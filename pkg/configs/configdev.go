package configs

const (
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

	// MYSQLTESTCONNSTRING - Mysql test connection string
	MYSQLTESTCONNSTRING = "root@tcp(127.0.0.1:3306)/clg_test_mysql"

	// TOKENLIFETIME - Jwt token lifetime
	TOKENLIFETIME = 5

	// TOKENSIGNKEY - Jwt token sign in key
	TOKENSIGNKEY = "secret"
)
