package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestGetSwiftCodeBranchUsingRealDb(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306" // db connection localhost:port
	cfg.DBName = "v1"           // dbname

	var db *sql.DB
	var err error

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

	handler := &DbHandler{DB: db}

	router := gin.Default()
	router.GET("v1/swift_codes/:swiftcode", func(c *gin.Context) {
		handler.GetDetailsOfSingleSwiftcode(c)
	})

	t.Run("Swifcode(headquarter) not empty", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/swift_codes/BCECCLRMXXX", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actual SwiftCodeData
		err = json.Unmarshal(w.Body.Bytes(), &actual)
		assert.NoError(t, err)

		// Assert headquarter properties
		assert.NotEmpty(t, actual.Address)
		assert.NotEmpty(t, actual.Name)
		assert.NotEmpty(t, actual.Countryiso2code)
		assert.NotEmpty(t, actual.Countryname)
		assert.True(t, actual.IsHeadquarter)
		assert.NotEmpty(t, actual.Swiftcode)

		// Assert branches
		assert.NotEmpty(t, actual.Branches, "Headquarter should have branches")

		if len(actual.Branches) > 0 {
			firstBranch := actual.Branches[0]
			assert.NotEmpty(t, firstBranch.Address)
			assert.NotEmpty(t, firstBranch.Name)
			assert.NotEmpty(t, firstBranch.Countryiso2code)
			assert.False(t, firstBranch.IsHeadquarter)
			assert.NotEmpty(t, firstBranch.Swiftcode)
		}
	})

	t.Run("Swiftcode(headquarter) equal value", func(t *testing.T) {
		// Use the known headquarters Swift code
		req, _ := http.NewRequest("GET", "/v1/swift_codes/BCECCLRMXXX", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actual SwiftCodeData
		err = json.Unmarshal(w.Body.Bytes(), &actual)
		assert.NoError(t, err)

		// Create the expected response structure
		expected := SwiftCodeData{
			Address:         "  ",
			Name:            "BANCO CENTRAL DE CHILE",
			Countryiso2code: "CL",
			Countryname:     "CHILE",
			IsHeadquarter:   true,
			Swiftcode:       "BCECCLRMXXX",
			Branches: []SwiftCodeData{
				{
					Address:         "  ",
					Name:            "BANCO CENTRAL DE CHILE",
					Countryiso2code: "CL",
					Countryname:     "CHILE",
					IsHeadquarter:   false,
					Swiftcode:       "BCECCLRMCSH",
				},
				{
					Address:         "  ",
					Name:            "BANCO CENTRAL DE CHILE",
					Countryiso2code: "CL",
					Countryname:     "CHILE",
					IsHeadquarter:   false,
					Swiftcode:       "BCECCLRMFCE",
				},
				{
					Address:         "  ",
					Name:            "BANCO CENTRAL DE CHILE",
					Countryiso2code: "CL",
					Countryname:     "CHILE",
					IsHeadquarter:   false,
					Swiftcode:       "BCECCLRMFES",
				},
				{
					Address:         "  ",
					Name:            "BANCO CENTRAL DE CHILE",
					Countryiso2code: "CL",
					Countryname:     "CHILE",
					IsHeadquarter:   false,
					Swiftcode:       "BCECCLRMFRP",
				},
			},
		}

		// Compare the entire structure
		assert.Equal(t, expected, actual)
	})

	t.Run("Swiftcode(branch)", func(t *testing.T) {

		req, _ := http.NewRequest("GET", "/v1/swift_codes/TPEOPLPWAAS", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actual SwiftCodeData
		err = json.Unmarshal(w.Body.Bytes(), &actual)
		assert.NoError(t, err)

		expected := SwiftCodeData{
			Address:         "FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",
			Name:            "PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA",
			Countryiso2code: "PL",
			Countryname:     "POLAND",
			Swiftcode:       "TPEOPLPWAAS",
			IsHeadquarter:   false,
		}
		assert.Equal(t, expected, actual)
	})

}

func TestReturnAllSwiftCodesByCountry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306" // db connection localhost:port
	cfg.DBName = "v1"           // dbname

	var db *sql.DB
	var err error

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

	handler := &DbHandler{DB: db}

	router := gin.Default()
	router.GET("v1/swift_codes/country/:countryiso2code", func(c *gin.Context) {
		handler.ReturnAllSwiftCodesCountry(c)
	})

	t.Run("Swiftcodes by countryISO2code", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/swift_codes/country/PL", nil) // sample
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actual_country CountryResponse
		err = json.Unmarshal(w.Body.Bytes(), &actual_country)
		assert.NoError(t, err)

		// Assert retreive properties
		assert.NotEmpty(t, actual_country.Countryiso2code)
		assert.NotEmpty(t, actual_country.Countryname)

		// Assert Swiftcodes
		assert.NotEmpty(t, actual_country.Swiftcodes, "List of swiftcodes")

		if len(actual_country.Swiftcodes) > 0 {
			first := actual_country.Swiftcodes[0]
			if first.IsHeadquarter == false {
				assert.NotEmpty(t, first.Address)
				assert.NotEmpty(t, first.BankName)
				assert.NotEmpty(t, first.Countryios2code)
				assert.False(t, first.IsHeadquarter)
				assert.NotEmpty(t, first.SwiftCode)

			} else {
				assert.NotEmpty(t, first.Address)
				assert.NotEmpty(t, first.BankName)
				assert.NotEmpty(t, first.Countryios2code)
				assert.True(t, first.IsHeadquarter)
				assert.NotEmpty(t, first.SwiftCode)
			}
		}
	})

}

func TestGetInfoSwiftCodeUsingMockDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		swiftcode        string
		mockSetup        func(mock sqlmock.Sqlmock)
		expectedStatus   int
		expectedResponse SwiftCodeData
	}{
		{
			name:      "Query for headquarter and brances used by swiftcode",
			swiftcode: "TPEOPLPWAAS",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"address", "name", "countryiso2code", "countryname", "swiftcode",
				}).AddRow(
					"FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",
					"PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA",
					"PL",
					"POLAND",
					"TPEOPLPWAAS",
				)
				mock.ExpectQuery("SELECT (.+) FROM swift_codes").WillReturnRows(rows)
			},

			expectedStatus: http.StatusOK,
			expectedResponse: SwiftCodeData{
				Address:         "FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",
				Name:            "PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA",
				Countryiso2code: "PL",
				Countryname:     "POLAND",
				Swiftcode:       "TPEOPLPWAAS",
				IsHeadquarter:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			handler := &DbHandler{DB: db}

			router := gin.New()
			router.GET("/v1/swift_codes/:swiftcode", func(c *gin.Context) {
				handler.GetDetailsOfSingleSwiftcode(c)
			})

			req, _ := http.NewRequest("GET", "/v1/swift_codes/"+tt.swiftcode, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var actual SwiftCodeData
			err = json.Unmarshal(w.Body.Bytes(), &actual)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedResponse, actual)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}

}
