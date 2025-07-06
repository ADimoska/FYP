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
	"time"

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

// ================================== Mountain ==================================

func getGGMultiplier(communityID int, currentTime string) float64 {
	c1GGtimes := map[string]bool{
		"0:30": true, "1:00": true, "1:30": true, "2:00": true, "2:30": true, "3:00": true,
		"3:30": true, "4:00": true, "4:30": true, "5:00": true, "5:30": true, "6:00": true,
		"6:30": true, "7:00": true, "7:30": true, "8:00": true, "8:30": true, "9:00": true,
		"9:30": true, "10:00": true, "10:30": true, "11:00": true, "11:30": true, "12:00": true,
		"12:30": true,
	}
	c2GGtimes := map[string]bool{
		"13:00": true, "13:30": true, "14:00": true, "14:30": true, "15:00": true, "15:30": true,
		"16:00": true, "16:30": true, "17:00": true, "17:30": true, "18:00": true, "18:30": true,
		"19:00": true, "19:30": true, "20:00": true, "20:30": true, "21:00": true, "21:30": true,
		"22:00": true, "22:30": true, "23:00": true, "23:30": true, "0:00": true,
	}

	switch communityID {
	case 1:
		if c1GGtimes[currentTime] {
			return 20.0
		}
	case 2:
		if c2GGtimes[currentTime] {
			return 20.0
		}
	}
	return 0.0
}

func GetSolarExposure(city string, dateStr string, dataStore map[string]map[int]map[string]map[int]float64) (float64, error) {
	// Parse the date string
	const layout = "2/01/2006"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %v", err)
	}

	// Extract components
	day := t.Day()
	year := t.Year()
	month := t.Month().String() // e.g., "January"

	// Convert to title case abbreviation to match map keys (e.g., "Jan", "Feb", ...)
	monthAbbr := month[:3]

	// Capitalize first letter and lowercase the rest
	monthAbbr = strings.Title(strings.ToLower(monthAbbr))

	// Retrieve value from the nested map
	if yearData, ok := dataStore[city]; ok {
		if monthData, ok := yearData[year]; ok {
			if dayData, ok := monthData[monthAbbr]; ok {
				if value, ok := dayData[day]; ok {
					return value, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("data not found for the given city and date")
}