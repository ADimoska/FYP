package main

import (
	"CES/house"
	"fmt"
	"os"
	"strconv"
	"encoding/csv"
	"time"
)

// func main() {
// 	myHouse := house.NewHouse("Ana", "London", 5.5 ) 
// 	fmt.Println("Customer:", myHouse.GetCustomer())
// 	myHouse.SetCustomer("Yomna")
// 	fmt.Println("Customer:", myHouse.GetCustomer())
//     fmt.Println("Location:", myHouse.GetLocation())
//     fmt.Printf("Generator Capacity: %.2f kW\n", myHouse.GetGen_capacity())
// }

func ReadCSVFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
	}
	
	return records, nil
}

// ParseHousesFromCSVRecords parses a 2D slice of CSV records and returns a slice of House pointers.
// It skips the header row and processes each subsequent row by extracting house ID, location, and capacity.
// Only unique house IDs (based on the first column) are added to the result.
// If capacity parsing fails for a row, that row is skipped with an error message.
func ParseHousesFromCSVRecords(records [][]string) []*house.House {
	var houses []*house.House
	addedHouse := "0"

	// Skip header and process each row
	for i, row := range records {
		if i == 0 {
			continue // skip header
		}

		capacity, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Printf("Error parsing capacity on line %d: %v\n", i+1, err)
			continue
		}

		if row[0] != addedHouse {
			h := house.NewHouse(row[0], row[2], capacity, 0.0, 0.0, 0.0, 0.0)
			houses = append(houses, h)
			addedHouse = row[0]
		}
	}

	return houses
}

// extractDates extracts the start and end dates from a 2D slice of CSV records.
// It skips the header row and looks for specific rows where the first column is "1" or "2".
// The start date is taken from the first row with "1" in the first column.
// The end date is set based on the most recent date before encountering a row with "2" in the first column.
// Returns the extracted start and end dates as strings. 
func extractDates(records [][]string) (string, string) {
	start_date_set := false
	end_date_set := false
	var start_date string
	var end_date string
	var most_recent_end_date string

	for i, row := range records {
		// Skip header and process each row
		if i == 0 {
			continue
		}

		if row[0] == "1" && !start_date_set {
			start_date = row[4]
			start_date_set = true
		} 

		if row[0] == "2" && !end_date_set {
			end_date = most_recent_end_date
			end_date_set = true
		} 

		most_recent_end_date = row[4]
	}

	return start_date, end_date
}

// iterateDates takes two date strings (start and end) in the format "2/01/2006",
// parses them into time.Time values, and prints each date from the start to the end, inclusive.
// It validates that both dates are correctly formatted and that the start date is not after the end date.
// If any validation fails, it prints an error message and exits the function.
func iterateDates(startDateStr, endDateStr string) {
	const layout = "2/01/2006"

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		fmt.Println("Invalid start date:", err)
		return
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		fmt.Println("Invalid end date:", err)
		return
	}

	// Ensure the start date is before or equal to end date
	if startDate.After(endDate) {
		fmt.Println("Start date is after end date")
		return
	}

	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		fmt.Println(currentDate.Format(layout))
	}
}

func main() {

	records, err := ReadCSVFile("2012_2013_Solar_home_electricity_data_v2.csv")
	if err != nil {
		fmt.Println("Failed to read CSV:", err)
		return
	}
	
	houses := ParseHousesFromCSVRecords(records)
	
	for i, h := range houses{
		fmt.Printf("House %d: %+v\n", i+1, *h)
	}

	

	start_date, end_date := extractDates(records)
	fmt.Printf("Final Start date: %s\n", start_date)
	fmt.Printf("Final End date: %s\n", end_date)
	//-> itterate over half hours

	iterateDates(start_date, end_date)

	//next step Iterate over periods, read the used energy and generated energy 

	// for i , h := range houses {
	// 	fmt.Printf("House %d: {customer:%s location:%s gen_capacity:%.2f}\n", i+1,  h.GetCustomer(), h.GetLocation(), h.GetGen_capacity())
	// }
	
}



