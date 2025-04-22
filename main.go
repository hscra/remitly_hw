package main

import (
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

	router.GET("/v1/swift_codes/:swiftcode", func(c *gin.Context) {
		swiftcodeHandler.GetDetailsOfSingleSwiftcode(c)
	})
	router.GET("/v1/swift_codes/country/:countryiso2code", swiftcodeHandler.ReturnAllSwiftCodesCountry)
	router.POST("/v1/:swift_codes", swiftcodeHandler.AddSwiftCodeToCountry)
	router.DELETE("/v1/swift_codes/:swiftcode", swiftcodeHandler.DeleteSwiftCode)

	router.Run("localhost:8080")

}
