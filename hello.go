package main

import (
	"context"
	"fmt"
	"io"

	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql" // for mysql driver
	"github.com/gocarina/gocsv"      // for unmarshall
)

type SwiftcodeData struct {
	CountryIso2Code string `csv:"COUNTRY ISO2 CODE"`
	SwiftCode string `csv:"SWIFT CODE"`
	CodeType string `csv:"CODE TYPE"`
	Name string `csv:"NAME"`
	Address string `csv:"ADDRESS"`
	TownName string `csv:"TOWN NAME"`
	CountryName string `csv:"COUNTRY NAME"`
	TimeZone string `csv:"TIME ZONE"`
}

type Location struct{
	Id int64
	lat float64
	lng float64 
}



var db *sql.DB

func main(){
	// fmt.Println("Heollo,world")
	// fmt.Println(quote.Go())
	// message := greetings.Hello("Glayds")
	// fmt.Println(message)

	// open file
	f, err := os.Open("swiftcode.csv")
	if err !=nil{
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

		fmt.Printf("%+v\n",rec)
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
	if err := gocsv.UnmarshalFile(f,&swiftcodes); err != nil{
		panic(err)
	}

	var lineExistMap = make(map[string]bool)
	// var duplicate string 
	// Print out the parsed data
	for index,swiftcode := range swiftcodes{
		// duplicate = r[0]
		if _, exist := lineExistMap[swiftcode.SwiftCode]; exist{
			continue
		} else {
			fmt.Println(index, swiftcode.SwiftCode)

		}
		// fmt.Println(swiftcode.CountryIso2Code)
		// fmt.Println(swiftcode.SwiftCode)
	}
	// Reset the file reader
	if _,err := f.Seek(0,0); err != nil{
		panic(err)
	}

	now := time.Now()
	readChannel := make(chan SwiftcodeData, 1)
	readFilePath := "swiftcode.csv"
	readFile, err := os.OpenFile(readFilePath,os.O_RDONLY, os.ModePerm)
	if err !=nil{
		panic(err)
	}
	defer readFile.Close()

	readFromCSV(readFile,readChannel)

	cnt :=0

	for r:= range readChannel{
		fmt.Println(r)
		cnt++
	}
	fmt.Println(time.Since(now),cnt)

	// import mysql driver
	// TODO make function for db connection

	// Capture connection properties

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")  //- if you have
	cfg.Net ="tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "v1"	// dbname 

	// get a database handle.
	// var err error
	db,err = sql.Open("mysql",cfg.FormatDSN())
	if err != nil{
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}
	fmt.Printf(("Connected!\n"))


	// db, err := sql.Open("mysql","user:root@/v1")
	// if err != nil{
	// 	panic(err)

	// }
	db.SetConnMaxLifetime(time.Minute *3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// check DB connections
	// locations, err := getAlllocations(1)
	// if err!= nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("locations %v\n",locations)


	createSwiftCodesTable(db)

}

func createSwiftCodesTable(db *sql.DB){
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
	_, err := db.ExecContext(ctx,createTableQuery)
	if err != nil{
		log.Fatalf("Error creating table :%v",err)
	}

	fmt.Println("Table 'users' created successfully (if it didn't exist).")
}



func getAlllocations(id float64)([]Location,error){
	var locations []Location

	rows, err := db.Query("SELECT * FROM locations")
	if err != nil {
		return nil, fmt.Errorf("rows %q,error %v",id,err)
	}
	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields
	for rows.Next(){
		var lc Location
		if err:= rows.Scan(&lc.Id,&lc.lat,&lc.lng); err != nil{
			return nil, fmt.Errorf("rows %q, error %v",id,err)

		}
		locations= append(locations, lc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows %q , error %v", id, err)
	}
	return locations, nil 
}



func readFromCSV(file *os.File, c chan SwiftcodeData ){
	// Set a pipe as the demiliter for reading
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader{
		r := csv.NewReader(in)
		r.Comma = ','
		return r
	})

	go func(){
		err := gocsv.UnmarshalToChan(file,c)
		if err != nil{
			panic(err)
		}
	}()


}

