package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsingCSV(t *testing.T) {
	t.Run("Basic parsing", func(t *testing.T) {
		tempDir := t.TempDir()
		testCSVPath := filepath.Join(tempDir, "test_swiftcodes.csv")

		// sample
		csvContent := `COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
AL,AAISALTRXXX,BIC11,UNITED BANK OF ALBANIA SH.A,"HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023",TIRANA,ALBANIA,Europe/Tirane
BG,ABIEBGS1XXX,BIC11,ABV INVESTMENTS LTD,"TSAR ASEN 20  VARNA, VARNA, 9002",VARNA,BULGARIA,Europe/Sofia
BG,ADCRBGS1XXX,BIC11,ADAMANT CAPITAL PARTNERS AD,"JAMES BOURCHIER BLVD 76A HILL TOWER SOFIA, SOFIA, 1421",SOFIA,BULGARIA,Europe/Sofia`

		// Write testCSV
		err := os.WriteFile(testCSVPath, []byte(csvContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test CSV file: %v", err)
		}

		file, err := os.Open(testCSVPath)
		if err != nil {
			t.Fatalf("Failed to open test CSV file: %v", err)
		}

		defer file.Close()

		c := make(chan SwiftcodeData, 10) // buffer for larger

		go func() {
			err := ReadFromCSV(file, c)
			if err != nil {
				t.Errorf("ReadFromCSV() returned error : %v", err)
			}
		}()

		// Gather results inoto []swiftcodes
		var swiftcodes []SwiftcodeData
		for data := range c {
			swiftcodes = append(swiftcodes, data)
		}

		// Assertations
		if len(swiftcodes) != 3 {
			t.Errorf("Expected 3 swift codes regardless of headquarter, got %d", len(swiftcodes))
		}

		if swiftcodes[0].CountryIso2Code != "AL" {
			t.Errorf("Exprected country ISO2 code 'AL', got '%v'", swiftcodes[0].CountryIso2Code)
		}

		if swiftcodes[1].SwiftCode != "ABIEBGS1XXX" {
			t.Errorf("Expected swift code 'ABIEBGS1XXX', got '%s'", swiftcodes[1].SwiftCode)
		}
	})

}
