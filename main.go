package main

import (
	"log"

	"swiftcode/database"
	"swiftcode/handlers"

	"github.com/gin-gonic/gin" // gin framework
)

func main() {
	// Initiation channel
	readChannel := database.Init()

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
	router.GET("/v1/swift_codes/country/:countryiso2code", func(c *gin.Context) {
		swiftcodeHandler.ReturnAllSwiftCodesCountry(c)
	})

	router.POST("/v1/:swift_codes", swiftcodeHandler.AddSwiftCodeToCountry)
	router.DELETE("/v1/swift_codes/:swiftcode", swiftcodeHandler.DeleteSwiftCode)

	router.Run(":8080")

}
