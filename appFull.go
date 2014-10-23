package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/hoisie/web"
)

//Here, I define the ype that will hold the data for each patient.
//One predicion "Y/N" and one probability associated to it.
type Record struct {
	Pred bool    `json:"pred"`
	P    float64 `json:"p"`
}

//The whole data structure
var data = map[string]Record{}

//Open the .csv file, load the data in the data structure
func ParseCsv(filename string) {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	reader := csv.NewReader(fi)
	reader.Comma = ',' //Just to remember that we can change these.

	lineCount := 0
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		lineCount += 1
		pred, _ := strconv.ParseBool(record[1])
		p, _ := strconv.ParseFloat(record[2], 64)
		data[record[0]] = Record{
			pred,
			p,
		}
	}

	fmt.Println("Read", lineCount, "Records.")

	fmt.Println("See for yourself", data)

}

//This is the handler for the request. It is the main code bloc, that does the work when a request is received.
func get(ctx *web.Context, val string) string {
	// 1. Parse
	epicID := ctx.Params["epicid"]
	record := data[epicID]
	fmt.Println(record)
	if record == (Record{}) {
		ctx.NotFound("This epicID doesn't match anything we have...sorry!")
		return ""
	}
	//Create the Json
	jsonResponse, _ := json.Marshal(record)
	//Set the contents
	ctx.SetHeader("X-Powered-By", "web.go", true)
	ctx.SetHeader("Connection", "close", true)
	ctx.SetHeader("Content-Type", "application/json", true)
	return string(jsonResponse)
}

//Main method
func main() {

	//Load the data
	ParseCsv("data.csv")

	//Launch the webserver
	web.Get("/prediction(.*)", get)
	web.Run("0.0.0.0:80")
}
