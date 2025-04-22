package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"swiftcode/database"
	"swiftcode/handlers"

	"github.com/gin-gonic/gin" // gin framework
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

func main() {
	// Read file to parsing csv
	now := time.Now()
	readChannel := make(chan database.SwiftcodeData, 1)
	readFilePath := "swiftcode.csv"
	readFile, err := os.OpenFile(readFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to open file at path %s: %v", readFilePath, err)
	}
	defer readFile.Close()
	database.ReadFromCSV(readFile, readChannel)

	cnt := 0
	fmt.Println(time.Since(now), cnt)

	// Database connection
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Printf("Error %s when connecting database", err)
	}
	defer db.Close()

	// Create swift_cdoes table
	database.CreateSwiftCodesTable(db)

	// Insert parsed csv into database
	for swiftcode := range readChannel {
		database.InsertSwiftCodes(db, swiftcode)
		if err != nil {
			log.Printf("Error inserting swift code: %v", err)
		}
	}

	// Seperate database and handler function.
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
