package main

import (
	"github.com/dcaraxes/gotenberg-go-client/v8"
	"log"
	"net/http"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	index, err := f.NewSheet("ExampleSheet")
	if err != nil {
		log.Fatal(err)
	}

	err = f.SetCellValue("ExampleSheet", "A2", "Hello world.")
	if err != nil {
		log.Fatal(err)
	}
	err = f.SetCellValue("ExampleSheet", "B2", 100)
	if err != nil {
		log.Fatal(err)
	}

	err = f.SetColWidth("ExampleSheet", "A", "B", 20)
	if err != nil {
		log.Fatal(err)
	}
	f.SetActiveSheet(index)
	if err := f.SaveAs("demo.xlsx"); err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	client := &gotenberg.Client{Hostname: "http://localhost:3000", HTTPClient: httpClient}

	xlsFile, err := gotenberg.NewDocumentFromPath("demo.xlsx", "demo.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	req := gotenberg.NewOfficeRequest(xlsFile)
	err = client.Store(req, "demo.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// Multiple files

	xlsFile2, err := gotenberg.NewDocumentFromPath("demo2.xlsx", "demo.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	req = gotenberg.NewOfficeRequest(xlsFile, xlsFile2)
	err = client.Store(req, "demo.zip")
	if err != nil {
		log.Fatal(err)
	}

	// Merge files
	req = gotenberg.NewOfficeRequest(xlsFile, xlsFile2)
	req.Merge()
	err = client.Store(req, "demo_merged.pdf")
	if err != nil {
		log.Fatal(err)
	}

}
