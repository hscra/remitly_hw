package database

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"time"

	"database/sql" // mysql API

	"github.com/go-sql-driver/mysql" // for mysql driver
	"github.com/gocarina/gocsv"      // for unmarshall
)

type SwiftcodeData struct {
	CountryIso2Code string `csv:"COUNTRY ISO2 CODE"`
	SwiftCode       string `csv:"SWIFT CODE"`
	CodeType        string `csv:"CODE TYPE"`
	Name            string `csv:"NAME"`
	Address         string `csv:"ADDRESS"`
	TownName        string `csv:"TOWN NAME"`
	CountryName     string `csv:"COUNTRY NAME"`
	TimeZone        string `csv:"TIME ZONE"`
}

// var db *sql.DB

func ReadFromCSV(file *os.File, c chan SwiftcodeData) error {
	// Set a pipe as the demiliter for reading
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = ','

		return r
	})

	go func() {
		err := gocsv.UnmarshalToChan(file, c)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func ConnectDatabase() (db *sql.DB, err error) {
	// Capture connection properties
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306" // db connection address:port
	cfg.DBName = "v1"           // dbname

	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("Failed to open your mysql database. Please check your environment set of DBUSER, DBPASS")
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println("Failed to ping() check for your db connection")
		panic(pingErr)

	}
	fmt.Printf(("Database connected!\n"))

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}

func CreateSwiftCodesTable(db *sql.DB) {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS swift_codes(
			id INT AUTO_INCREMENT PRIMARY KEY,
			countryiso2code VARCHAR(255) NOT NULL,
			swiftcode VARCHAR(255) UNIQUE NOT NULL,
			codetype VARCHAR(255),
			name VARCHAR(255),
			address VARCHAR(255),
			townname VARCHAR(255),
			countryname VARCHAR(255),
			timezone VARCHAR(255)
		);
	`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Fatalf("Error creating table :%v", err)
	}

	fmt.Println("Table 'swift_codes' created successfully (if it didn't exist).")
}

func InsertSwiftCodes(db *sql.DB, sc SwiftcodeData) {
	query := `
	INSERT IGNORE INTO swift_codes (
		countryiso2code,
		swiftcode,
		codetype,
		name,
		address,
		townname,
		countryname,
		timezone
		)
		VALUES (?,?,?,?,?,?,?,?)
	`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalf("Error %s when preparing SQL query\n", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(
		ctx,
		sc.CountryIso2Code,
		sc.SwiftCode,
		sc.CodeType,
		sc.Name,
		sc.Address,
		sc.TownName,
		sc.CountryName,
		sc.TimeZone,
	)
	if err != nil {
		log.Fatalf("Error %s when inserting row into table\n", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Fatalf("Error %s when finding rows affected\n", err)
	}
	// log.Println("%d swiftcodes created\n")
}
