module example/hello

go 1.22.5

replace example.com/greetings => ../greetings

require (
	github.com/go-sql-driver/mysql v1.9.2
	github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1
)

require filippo.io/edwards25519 v1.1.0 // indirect
