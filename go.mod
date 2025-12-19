module student-portal

go 1.25.1

require (
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.46.0
	golang.org/x/oauth2 v0.34.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect
