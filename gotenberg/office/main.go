package main

import (
	"context"
	"github.com/starwalkn/gotenberg-go-client/v8"
	"github.com/starwalkn/gotenberg-go-client/v8/document"
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
	client, err := gotenberg.NewClient("http://localhost:3000", httpClient)
	if err != nil {
		log.Fatal(err)
	}

	xlsFile, err := document.FromPath("demo.xlsx", "demo.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	req := gotenberg.NewLibreOfficeRequest(xlsFile)
	err = client.Store(ctx, req, "demo.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// Multiple files

	xlsFile2, err := document.FromPath("demo2.xlsx", "demo.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	req = gotenberg.NewLibreOfficeRequest(xlsFile, xlsFile2)
	err = client.Store(ctx, req, "demo.zip")
	if err != nil {
		log.Fatal(err)
	}

	// Merge files
	req = gotenberg.NewLibreOfficeRequest(xlsFile, xlsFile2)
	req.Merge()
	err = client.Store(ctx, req, "demo_merged.pdf")
	if err != nil {
		log.Fatal(err)
	}

}
