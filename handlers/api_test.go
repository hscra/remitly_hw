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

func TestGetSwiftCodeUsingRealDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306" // db connection address:port
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
	router.GET("/swiftcode/:swiftcode", func(c *gin.Context) {
		handler.GetDetailsOfSingleSwiftcode(c)
	})

	req, _ := http.NewRequest("GET", "/swiftcode/TPEOPLPWAAS", nil)
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
			router.GET("/swiftcode/:swiftcode", func(c *gin.Context) {
				handler.GetDetailsOfSingleSwiftcode(c)
			})

			req, _ := http.NewRequest("GET", "/swiftcode/"+tt.swiftcode, nil)
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
