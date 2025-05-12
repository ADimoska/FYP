package main

import (
	"CES/house"
	"fmt"
	"os"
	"strconv"
	
	"encoding/csv"
	"time"
	"sort"
	
	

)


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
func iterateDates(startDateStr, endDateStr string, houses []*house.House) {
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
		date_str := currentDate.Format(layout)
		fmt.Println(date_str)
		executeDate(date_str, houses)
	}
}

func executeDate (date string, houses []*house.House) {
	timeArray := []string{"0:30","1:00","1:30","2:00","2:30","3:00","3:30","4:00","4:30","5:00","5:30","6:00","6:30","7:00","7:30","8:00","8:30","9:00","9:30","10:00","10:30","11:00","11:30","12:00","12:30","13:00","13:30","14:00","14:30","15:00","15:30","16:00","16:30","17:00","17:30","18:00","18:30","19:00","19:30","20:00","20:30","21:00","21:30","22:00","22:30","23:00","23:30","0:00"}
	for _, time := range timeArray {
		executeTime(date, time, houses)

		//used just for testing 
		// if date == "26/04/2013" && time == "0:30" {
		// 	for _, h := range houses {
		// 		println(h.GetCL(), h.GetGC(), h.GetGG())
		// 	}
		// }
	}
}

func executeTime(date, time string, houses []*house.House) {
	for _, h := range houses {
		h.GetCurrentEnergy(date, time)

	}
}

func processData(records [][]string, houses []*house.House) {
	const layout = "2/01/2006"

	// Skip header
	for i, row := range records {
		if i == 0 {
			continue
		}
		if i == 1  {
			continue
		}

		// CSV Structure: Customer, Capacity, Postcode, Type, Date, 0:30, 1:00, ..., 0:00
		customer := row[0]
		consumptionType := row[3]
		date := row[4]

		// Collect time data starting from index 5
		for periodIndex := 5; periodIndex < len(row)-1; periodIndex++ {
			periodTime := records[1][periodIndex] // Header has the time labels

			valueStr := row[periodIndex]
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				fmt.Printf("Skipping invalid value at row %d, period %d: %v\n", i+1, periodIndex, err)
				continue
			}

			// Find the matching house
			for _, h := range houses {
				if h.GetCustomer() == customer {
					h.StoreEnergyData(date, periodTime, consumptionType, value)
					break
				}
			}
		}
	}
}


func SaveEnergyDataToCSV(h *house.House, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Date", "Time", "ConsumptionType", "Value"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	data := h.GetEnergyData()

	// Sort date keys chronologically
	type dateEntry struct {
		str  string
		time time.Time
	}
	const dateLayout = "2/01/2006"

	var dateList []dateEntry
	for dateStr := range data {
		parsedDate, err := time.Parse(dateLayout, dateStr)
		if err != nil {
			fmt.Printf("Skipping invalid date %s: %v\n", dateStr, err)
			continue
		}
		dateList = append(dateList, dateEntry{str: dateStr, time: parsedDate})
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i].time.Before(dateList[j].time)
	})

	// Process each sorted date
	for _, date := range dateList {
		times := data[date.str]

		// Sort time keys using actual time
		type timeEntry struct {
			str  string
			time time.Time
		}
		const timeLayout = "15:04" // 24-hour clock

		var timeList []timeEntry
		for timeStr := range times {
			parsedTime, err := time.Parse(timeLayout, timeStr)
			if err != nil {
				fmt.Printf("Skipping invalid time %s: %v\n", timeStr, err)
				continue
			}
			timeList = append(timeList, timeEntry{str: timeStr, time: parsedTime})
		}
		sort.Slice(timeList, func(i, j int) bool {
			return timeList[i].time.Before(timeList[j].time)
		})

		// For each sorted time
		for _, timeItem := range timeList {
			types := times[timeItem.str]

			// Sort consumption types
			ctypeKeys := make([]string, 0, len(types))
			for ctype := range types {
				ctypeKeys = append(ctypeKeys, ctype)
			}
			sort.Strings(ctypeKeys)

			// Write the records
			for _, ctype := range ctypeKeys {
				value := types[ctype]
				record := []string{date.str, timeItem.str, ctype, fmt.Sprintf("%.4f", value)}
				if err := writer.Write(record); err != nil {
					return fmt.Errorf("failed to write record: %w", err)
				}
			}
		}
	}

	fmt.Println("Data written to", filename)
	return nil
}



func removeHousesUnder365(houses []*house.House) []*house.House {
	var filteredHouses []*house.House

	for _, h := range houses {
		energyData := h.GetEnergyData()

		// Check first level: number of dates
		if len(energyData) < 365 {
			continue
		}

		valid := true

		for _, times := range energyData {
			for _, consumptionTypes := range times {
				if len(consumptionTypes) != 3 {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}

		if valid {
			filteredHouses = append(filteredHouses, h)
		}
	}

	return filteredHouses
}

func removeHousesOtherCity(houses []*house.House) []*house.House {
	var filteredHouses []*house.House

	for _, h := range houses {
		city := h.GetCity()

		// Check first level: number of dates
		if city == "Other" {
			continue
		}

		// valid := true


		// if valid {
		filteredHouses = append(filteredHouses, h)
		// }
	}

	return filteredHouses
}

func classifyCity(lat, lng float64) string {
	switch {
	case lat >= -34.2 && lat <= -33.6 && lng >= 150.8 && lng <= 151.4:
		return "Sydney"
	case lat >= -33.6 && lat <= -33.1 && lng >= 151.2 && lng <= 151.5:
		return "Gosford"
	case lat >= -33.2 && lat <= -32.7 && lng >= 151.5 && lng <= 151.9:
		return "Newcastle"
	case lat >= -33.3 && lat <= -32.6 && lng >= 150.9 && lng <= 151.5:
		return "Cessnock"
	default:
		return "Other"
	}
}

func ReadPostcodeLatLngMap(filePath string) (map[string][]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %v", err)
	}

	postcodeLatLng := make(map[string][]float64)

	for i, row := range records {
		// Skip header
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}
		postcode := row[0]
		lat, err1 := strconv.ParseFloat(row[1], 64)
		lng, err2 := strconv.ParseFloat(row[2], 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid lat/lng at row %d: %v, %v", i, err1, err2)
		}
		postcodeLatLng[postcode] = []float64{lat, lng}
	}

	return postcodeLatLng, nil
}






func SetCityToHouses (map_loc map[string][]float64, houses []*house.House){
	for _, h := range houses{
		loc := h.GetLocation()
		lat_lng := map_loc[loc]
		city := classifyCity(lat_lng[0], lat_lng[1])
		h.SetCity(city)
		//fmt.Printf("House %s: %s,%s\n", h.GetCustomer(), city, h.GetLocation())
	}
}
func main() {

	records, err := ReadCSVFile("2012_2013_Solar_home_electricity_data_v2.csv")
	if err != nil {
		fmt.Println("Failed to read CSV:", err)
		return
	}

	map_loc, err := ReadPostcodeLatLngMap("postcodes.csv")
	if err != nil {
		fmt.Println("Failed to read postcodes.csv:", err)
		return
	}
	
	houses := ParseHousesFromCSVRecords(records)
	
	for i, h := range houses{
		fmt.Printf("House %d: %+v\n", i+1, *h)
	}

	

	start_date, end_date := extractDates(records)
	fmt.Printf("Final Start date: %s\n", start_date)
	fmt.Printf("Final End date: %s\n", end_date)
	

	processData(records, houses)


	houses = removeHousesUnder365(houses) 

	SetCityToHouses(map_loc, houses)

	houses = removeHousesOtherCity(houses)

	for _, h := range houses{
		fmt.Printf("House %s:%s\n", h.GetCustomer(), h.GetCity())
	}

	iterateDates(start_date, end_date, houses)

	//used just for testing
	data := houses[0].GetEnergyData()
	fmt.Println(data["1/07/2012"]["0:30"]["GC"])
	
	// used just for testings
	for _, h := range houses {
		if h.GetCustomer() == "11" {
			err := SaveEnergyDataToCSV(h, "consumer_11_energydata.csv")
			if err != nil {
				fmt.Println("Error saving data:", err)
			}
			break
		}
	}
	
	// data_weather, err := readCSVToMonthDayMap("weather.csv")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// used just for testing
	// for day, value := range data_weather["Feb"] {
	// 	fmt.Printf("Feb %d: %.2f\n", day, value)
	// }
	
	dataStore := make(map[string]map[int]map[string]map[int]float64)
	loadAllCityYearData(dataStore)

    // Accessing an example value:
	fmt.Println("--")
    fmt.Println(dataStore["Sydney"][2012]["Jan"][1])
	fmt.Println(dataStore["Sydney"][2013]["Jan"][1])
	fmt.Println(dataStore["Gosford"][2012]["Jan"][1])
	fmt.Println(dataStore["Gosford"][2013]["Jan"][1])
	fmt.Println(dataStore["Newcastle"][2012]["Jan"][1])
	fmt.Println(dataStore["Newcastle"][2013]["Jan"][1])
	fmt.Println(dataStore["Cessnock"][2012]["Jan"][1])
	fmt.Println(dataStore["Cessnock"][2013]["Jan"][1])

	// fmt.Println(dataStore["Sydney"][2012].key())
	// fmt.Println(dataStore)
	// for key := range dataStore["Sydney"][2012] {
	// 	fmt.Println("*", key, "*")
	// 	fmt.Println(len(key)) 
	// }
	
		
}



