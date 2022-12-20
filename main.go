package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type downloadBody struct {
	URL string `json:"url"`
}

func main() {

	now := time.Now()
	aMonthAgo := now.AddDate(0, -1, 0)

	var year string
	flag.StringVar(&year, "year", aMonthAgo.Format("2006"), "the year to export the SAFT, defaults year of last month")

	var month string
	flag.StringVar(&month, "month", aMonthAgo.Format("01"), "the month to export the SAFT, defaults to last month")

	var destinationFolder string
	flag.StringVar(&destinationFolder, "destination", "", "the destination directory of the output file")

	var monthYearPattern string
	flag.StringVar(&monthYearPattern, "month-year-pattern", "%s/%s-%s", "pattern to append the year and month to the destination directory")

	flag.Parse()

	if year == "" {
		log.Fatal("year cannot be empty, use -year flag")
	}

	if month == "" {
		log.Fatal("month cannot be empty, use -month flag")
	}

	if monthYearPattern == "" {
		log.Fatal("month-year-pattern cannot be empty")
	}

	finalDestination := fmt.Sprintf(monthYearPattern, destinationFolder, year, month)

	accountName := os.Getenv("INVOICE_ACCOUNT_NAME")
	apiKey := os.Getenv("INVOICE_API_KEY")

	url := fmt.Sprintf("https://%s.app.invoicexpress.com/api/export_saft.json?month=%s&years=%s&api_key=%s",
		accountName,
		month,
		year,
		apiKey,
	)

	log.Printf("going to export SAFT for %s/%s into '%s'", year, month, finalDestination)

	success := false
	for !success {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("accept", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("error requesting SAFT: %s", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("error reading SAFT request response: %s", err)
		}
		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusAccepted:
			log.Println("waiting for SAFT...")
			time.Sleep(10 * time.Second)
		case http.StatusOK:
			downloadInfo := downloadBody{}
			json.Unmarshal(body, &downloadInfo)
			log.Printf("downloading SAFT from %s...", downloadInfo.URL)

			req, _ := http.NewRequest("GET", downloadInfo.URL, nil)
			req.Header.Add("accept", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatalf("error requestion SAFT download: %s", err)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatalf("error downloading SAFT: %s", err)
			}
			defer res.Body.Close()

			outputPath := fmt.Sprintf("%s/SAF-T_%s_%s.zip", finalDestination, month, year)
			err = os.WriteFile(outputPath, body, os.FileMode(0700))
			if err != nil {
				log.Fatalf("error writing SAFT to file: %s", err)
			}
			log.Printf("the SAFT file was written to %s", outputPath)
			success = true
		default:
			log.Fatalf("error reading SAFT request response: %s", res.Status)
		}
	}
}
