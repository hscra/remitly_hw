package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"swiftcode/database"
	"swiftcode/handlers"

	"github.com/gin-gonic/gin"  // gin framework
	"github.com/gocarina/gocsv" // for unmarshall
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

type Location struct {
	Id  int64
	lat float64
	lng float64
}

var db *sql.DB

func main() {
	// fmt.Println("Heollo,world")
	// fmt.Println(quote.Go())
	// message := greetings.Hello("Glayds")
	// fmt.Println(message)

	// open file
	f, err := os.Open("swiftcode.csv")
	if err != nil {
		log.Fatal(err)
	}

	// close the file
	defer f.Close()

	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", rec)
	}
	// Reset the file reader
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	// gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader{
	// 	r := csv.NewReader(in)
	// 	r.Comma = ' '
	// 	return r
	// })

	// To read csv file
	swiftcodes := []*SwiftcodeData{}
	if err := gocsv.UnmarshalFile(f, &swiftcodes); err != nil {
		panic(err)
	}

	// fmt.Printf("parsed? %+s\n",swiftcodes);

	/*Check duplicate and remove it*/
	var lineExistMap = make(map[string]bool)
	// var duplicate string
	// Print out the parsed data
	for index, swiftcode := range swiftcodes {
		// duplicate = r[0]
		if _, exist := lineExistMap[swiftcode.SwiftCode]; exist {
			continue
		} else {
			fmt.Println(index, swiftcode.SwiftCode)

		}
		// fmt.Println(swiftcode.CountryIso2Code)
		// fmt.Println(swiftcode.SwiftCode)
	}
	// Reset the file reader
	if _, err := f.Seek(0, 0); err != nil {
		panic(err)
	}

	// Another method to readCSV 확정 !!
	now := time.Now()
	readChannel := make(chan database.SwiftcodeData, 1)
	readFilePath := "swiftcode.csv"
	readFile, err := os.OpenFile(readFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer readFile.Close()
	database.ReadFromCSV(readFile, readChannel)

	// readFromCSV(readFile, readChannel)

	cnt := 0

	// for r:= range readChannel{
	// 	fmt.Println(r)
	// 	cnt++
	// }
	fmt.Println(time.Since(now), cnt)

	// import mysql driver
	// TODO make function for db connection

	db, err := database.ConnectDatabase()
	if err != nil {
		log.Printf("Error %s when connecting database", err)
	}
	defer db.Close()

	// Create table
	database.CreateSwiftCodesTable(db)

	// Insert parsed csv into database
	// for _,swiftcode := range swiftcodes {
	// 	insertSwiftCodes(db,*swiftcode)
	// }
	for swiftcode := range readChannel {
		database.InsertSwiftCodes(db, swiftcode)
		if err != nil {
			log.Printf("Error inserting swift code: %v", err)
		}
	}

	// db function, handler function seperate it.
	swiftcodeHandler := handlers.SwiftCodeHandler(db)
	router := gin.Default()

	router.GET("/v1/swift_codes/:swiftcode", swiftcodeHandler.GetDetailsOfSingleSwiftcode)
	router.GET("/v1/swift_codes/country/:countryiso2code", swiftcodeHandler.ReturnAllSwiftCodesCountry)

	router.Run("localhost:8080")

}

func insertSwiftCodes(db *sql.DB, sc SwiftcodeData) {
	query := `
	INSERT INTO swift_codes (
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
	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error %s when finding rows affected\n", err)
	}
	log.Printf("%d swiftcodes created\n", rows)
}

func createSwiftCodesTable(db *sql.DB) {
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

// func getAlllocations(id float64)([]Location,error){
// 	var locations []Location

// 	rows, err := db.Query("SELECT * FROM locations")
// 	if err != nil {
// 		return nil, fmt.Errorf("rows %q,error %v",id,err)
// 	}
// 	defer rows.Close()

// loop through rows, using Scan to assign column data to struct fields
// 	for rows.Next(){
// 		var lc Location
// 		if err:= rows.Scan(&lc.Id,&lc.lat,&lc.lng); err != nil{
// 			return nil, fmt.Errorf("rows %q, error %v",id,err)

// 		}
// 		locations= append(locations, lc)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("rows %q , error %v", id, err)
// 	}
// 	return locations, nil
// }

func readFromCSV(file *os.File, c chan SwiftcodeData) {
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

}
