//Weather helper functions
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"encoding/csv"
	"log"
	"regexp"

)

var nonLetters = regexp.MustCompile(`[^a-zA-Z]+`)

func cleanToLettersOnly(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.TrimSpace(s)
	return nonLetters.ReplaceAllString(s, "")
}

func readCSVToMonthDayMap(filename string) (map[string]map[int]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read the header to get month names
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	cleanHeaders := make([]string, len(header))
	monthData := make(map[string]map[int]float64)
	for i, month := range header {
		clean_month := cleanToLettersOnly(month)
		cleanHeaders[i] = clean_month
		monthData[clean_month] = make(map[int]float64)
	}

	// Read the data rows
	day := 1
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		for i, val := range record {
			if val == "" {
				continue // skip missing values
			}
			month := cleanHeaders[i]
			floatVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid float at day %d in month %s: %v", day, month, err)
			}
			monthData[month][day] = floatVal
		}

		day++
	}

	return monthData, nil
}

func storeCityYearWeatherData(city string, year int, filename string, dataStore map[string]map[int]map[string]map[int]float64) error {
    monthDayMap, err := readCSVToMonthDayMap(filename)
    if err != nil {
        return err
    }

    // Initialize nested maps if needed
    if _, ok := dataStore[city]; !ok {
        dataStore[city] = make(map[int]map[string]map[int]float64)
    }

    dataStore[city][year] = monthDayMap
    return nil
}

func loadAllCityYearData(dataStore map[string]map[int]map[string]map[int]float64) {
	cities := []string{"Sydney", "Gosford", "Newcastle", "Cessnock"}
	years := []int{2012, 2013}

	for _, city := range cities {
		for _, year := range years {
			fileName := fmt.Sprintf("solar_exposure/%s/Sol_expo_%s_%d_kwh_csv.csv", city, city, year)

			err := storeCityYearWeatherData(city, year, fileName, dataStore)
			if err != nil {
				log.Fatalf("Error loading data for %s %d: %v", city, year, err)
			}
		}
	}
}