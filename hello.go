package main

import (
	"fmt"
	"io"

	"encoding/csv"
	"log"
	"os"
)
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




}

