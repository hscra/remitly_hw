package main

import (
	"fmt"
	"io"

	"encoding/csv"
	"log"
	"os"

	"github.com/gocarina/gocsv"
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
	// Print out the parsed data
	for _,swiftcode := range swiftcodes{
		fmt.Println(swiftcode.CountryIso2Code)
	}
	// Reset the file reader
	if _,err := f.Seek(0,0); err != nil{
		panic(err)
	}

	



}

