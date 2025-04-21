package handlers

import (
	"database/sql"
	"fmt"
	"net/http" // for http.StatusOK like
	"strings"

	"github.com/gin-gonic/gin" // for gin web framework
)

type DbHandler struct {
	DB *sql.DB
}

type SwiftCodeData struct {
	Address         string          `json:"address"`
	Name            string          `json:"bankName"`
	Countryiso2code string          `json:"countryISO2"`
	Countryname     string          `json:"countryName"`
	IsHeadquarter   bool            `json:"isHeadquarter"`
	Swiftcode       string          `json:"siwftCode"`
	Branches        []SwiftCodeData `json:"branches,omitempty"`
}

func SwiftCodeHandler(db *sql.DB) *DbHandler {
	return &DbHandler{DB: db}
}

func (h *DbHandler) GetDetailsOfSingleSwiftcode(c *gin.Context) {
	fmt.Println("***REQUEST RECEIVED***")
	swiftcode := c.Param("swiftcode")

	var sc SwiftCodeData
	var branches []SwiftCodeData

	completeResponse := SwiftCodeData{
		IsHeadquarter: strings.Contains(swiftcode, "XXX"),
	}

	// Query row with matching swiftcode input
	row := h.DB.QueryRow("SELECT address, name, countryiso2code, countryname, swiftcode FROM swift_codes WHERE swiftcode=?", swiftcode)
	err := row.Scan(&completeResponse.Address, &completeResponse.Name, &completeResponse.Countryiso2code, &completeResponse.Countryname, &completeResponse.Swiftcode)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
			return
		}
		fmt.Printf("Database error for swiftcode %s: %v\n", swiftcode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// check swiftcode whether it indicates headquarter
	if strings.Contains(swiftcode, "XXX") {
		fmt.Println("This SWIFTCODE is for headquarter")
		sc.IsHeadquarter = true
	} else {
		fmt.Println("This SWIFTCODE is for branch")
		sc.IsHeadquarter = false
	}

	if sc.IsHeadquarter { // headquarter
		// find its branches
		brSwiftCode := swiftcode[:8]
		rows, err := h.DB.Query("SELECT address, name ,countryiso2code,countryname,swiftcode FROM swift_codes WHERE swiftcode LIKE ? and swiftcode != ?", brSwiftCode+"%", swiftcode)
		if err != nil {
			fmt.Printf("Error to query branches : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error to query"})
			return
		}
		defer rows.Close()

		// loop through rows, using Scan to assign column data to struct fields
		for rows.Next() {
			var branch SwiftCodeData
			if err := rows.Scan(&branch.Address, &branch.Name, &branch.Countryiso2code, &branch.Countryname, &branch.Swiftcode); err != nil {
				fmt.Printf("Error to qeury branches over the loop : %v", err)
				continue
			}
			branches = append(branches, branch)
		}
		completeResponse.Branches = branches

		if err := rows.Err(); err != nil {
			fmt.Printf("Error to scanning rows : %v", err)
		}

		c.JSON(http.StatusOK, completeResponse)

	} else { // branch

		c.JSON(http.StatusOK, gin.H{
			"address":       sc.Address,
			"bankName":      sc.Name,
			"countryISO2":   sc.Countryiso2code,
			"countryName":   sc.Countryname,
			"isHeadquarter": false,
			"siwftCode":     sc.Swiftcode,
		})

	}

}
