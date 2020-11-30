module alphaflow

go 1.15

require (
	github.com/attache/attache v0.7.0
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
)

replace github.com/attache/attache v0.7.0 => ../attache
