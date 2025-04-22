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

type SwiftCodeData struct { // for Endpoint 1,2,3,4
	Address         string          `json:"address"`
	Name            string          `json:"bankName"`
	Countryiso2code string          `json:"countryISO2"`
	Countryname     string          `json:"countryName"`
	IsHeadquarter   bool            `json:"isHeadquarter"`
	Swiftcode       string          `json:"siwftCode"`
	Branches        []SwiftCodeData `json:"branches,omitempty"`
}

type SwiftCodeSummary struct { // for Endpoint 2
	Address         string `json:"address"`
	BankName        string `json:"bankName"`
	Countryios2code string `json:"countryISO2"`
	IsHeadquarter   bool   `json:"isHeadquarter"`
	SwiftCode       string `json:"swiftCode"`
}

type CountryResponse struct { // for Endpoint 2
	Countryiso2code string             `json:"countryISO2"`
	Countryname     string             `json:"countryName"`
	Swiftcodes      []SwiftCodeSummary `json:"swiftCodes,omitempty"`
}

type NewSwiftCode struct { // for Endopoint POST 3
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func SwiftCodeHandler(db *sql.DB) *DbHandler {
	return &DbHandler{DB: db}
}

func (h *DbHandler) GetDetailsOfSingleSwiftcode(c *gin.Context) SwiftCodeData {
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
			return SwiftCodeData{}
		}
		fmt.Printf("Database error for swiftcode %s: %v\n", swiftcode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return SwiftCodeData{}
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
		return completeResponse

	} else { // branch

		c.JSON(http.StatusOK, completeResponse)
		return completeResponse

	}
}

func (h *DbHandler) ReturnAllSwiftCodesCountry(c *gin.Context) {
	fmt.Println("***REQUEST RECEIVED***")
	countryISO2code := c.Param("countryiso2code")

	var sc SwiftCodeSummary
	var swiftcodes []SwiftCodeSummary

	// Query row with matching countryISO2code
	rows, err := h.DB.Query("SELECT countryiso2code, countryname, swiftcode, address, name FROM swift_codes WHERE countryiso2code=?", countryISO2code)
	if err != nil {
		fmt.Printf("Database error for swift_codes :%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	defer rows.Close()

	countryResponse := CountryResponse{
		Countryiso2code: countryISO2code,
		Countryname:     "",
	}

	for rows.Next() {
		var countryName string

		if err := rows.Scan(&sc.Countryios2code, &countryName, &sc.SwiftCode, &sc.Address, &sc.BankName); err != nil {
			fmt.Printf("Error scanning row :%v", err)
			continue
		}

		if countryResponse.Countryname == "" {
			countryResponse.Countryname = countryName
		}

		sc.IsHeadquarter = false

		if strings.Contains(sc.SwiftCode, "XXX") {
			sc.IsHeadquarter = true
		}

		swiftcodes = append(swiftcodes, sc)

	}

	countryResponse.Swiftcodes = swiftcodes

	if err := rows.Err(); err != nil {
		fmt.Printf("Error to scanning rows : %v", err)
	}

	c.JSON(http.StatusOK, countryResponse)

}

func (h *DbHandler) AddSwiftCodeToCountry(c *gin.Context) {
	fmt.Println("***REQUEST RECEIVED***")
	var newSwiftCode NewSwiftCode

	if err := c.BindJSON(&newSwiftCode); err != nil {
		fmt.Printf("BindJSON error for swifr code")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Query row with matching swiftcode
	result, err := h.DB.Exec("INSERT INTO swift_codes(address,name,countryiso2code,countryname,swiftcode) VALUES (?,?,?,?,?) ",
		newSwiftCode.Address, newSwiftCode.BankName, newSwiftCode.CountryISO2, newSwiftCode.CountryName, newSwiftCode.SwiftCode)

	if err != nil {
		fmt.Printf("Database error for swift_codes :%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting rows affected: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not added"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully added SWIFT code: %s", newSwiftCode.SwiftCode),
	})

}

func (h *DbHandler) DeleteSwiftCode(c *gin.Context) {
	fmt.Println("***REQUEST RECEIVED***")
	swiftcode := c.Param("swiftcode")

	// Query row with matching swiftcode
	result, err := h.DB.Exec("DELETE FROM swift_codes WHERE swiftcode=?", swiftcode)
	if err != nil {
		fmt.Printf("Database error for swift_codes :%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting rows affected: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if rowsAffected == 0 {
		// No rows were deleted - the SWIFT code wasn't found
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted SWIFT code: %s", swiftcode),
	})

}
