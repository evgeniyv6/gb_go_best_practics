package main

import (
	"flag"
	"gb_go_best_practics/homework/sqlcsv/printer"
	"log"
	"os"
)

// go run main.go -request "select * from test1.csv where date>'2020-02-29' or iso_code='AFG'"
// go run main.go -request "select date, total_cases from test1.csv where date>'2020-02-29' and iso_code='AFG'"

var (
	request string
)

func init() {
	flag.StringVar(&request, "request", "", "sql select")
	flag.Parse()

	if request == "" {
		log.Println("sql request is empty")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	err := printer.Print(request)
	if err != nil {
		log.Fatal(err)
	}

}
