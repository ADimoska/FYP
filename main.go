package main

import (
	"CES/house"
	"CES/pool"
	"fmt"	
	"time"
	// "os"
	// "github.com/go-echarts/go-echarts/v2/charts"
	// "github.com/go-echarts/go-echarts/v2/opts"

	
)

var PoolBattery []float64
// iterateDates takes two date strings (start and end) in the format "2/01/2006",
// parses them into time.Time values, and prints each date from the start to the end, inclusive.
// It validates that both dates are correctly formatted and that the start date is not after the end date.
// If any validation fails, it prints an error message and exits the function.
func iterateDates(startDateStr, endDateStr string, houses []*house.House, p *pool.Pool) {
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
		executeDate(date_str, houses, p)
	}
	GenerateBlackoutChart(houses, "blackouts.html")
	
	CallPythonToGenerateHistogram(houses)
	CallPythonToPlotPoolBattery(PoolBattery)
	fmt.Println("Average threshold:",(total3DaySum/total3DayCount)) //E3 
	
}

var total3DaySum float64 = 0.0 //E3
var total3DayCount float64 = 0 //E3
func executeDate (date string, houses []*house.House, p *pool.Pool) {
	timeArray := []string{"0:30","1:00","1:30","2:00","2:30","3:00","3:30","4:00","4:30","5:00","5:30","6:00","6:30","7:00","7:30","8:00","8:30","9:00","9:30","10:00","10:30","11:00","11:30","12:00","12:30","13:00","13:30","14:00","14:30","15:00","15:30","16:00","16:30","17:00","17:30","18:00","18:30","19:00","19:30","20:00","20:30","21:00","21:30","22:00","22:30","23:00","23:30","0:00"}	
	// E3 start =================================================================
	if date == "1/07/2012" || date == "2/07/2012" || date == "3/07/2012"{
		for _, h := range houses {
		h.SetNext3Days(date)
		total3DaySum += h.Getlast3DaysConsumption()
		total3DayCount++
		}
	} else {
		for _, h := range houses {
			h.SetLast3Days(date)
			total3DaySum += h.Getlast3DaysConsumption()
			total3DayCount++
			}
	}
	// E3 end =================================================================
	
	for _, time := range timeArray {
		executeTime(date, time, houses, p)

		//used just for testing 
		// if date == "26/04/2013" && time == "0:30" {
		// 	for _, h := range houses {
		// 		println(h.GetCL(), h.GetGC(), h.GetGG())
		// 	}
		// }
		
	}
	// for _, h := range houses {
	// 	if h.GetCustomer() == "1" {
	// 		fmt.Println(h.GetBlackouts())
	// 	}
	// }
}

func executeTime(date, time string, houses []*house.House, p *pool.Pool) {
	for _, h := range houses {
		h.GetCurrentEnergy(date, time)
		h.AddBattery(-(h.GetGC()))
		h.AddBattery(-(h.GetCL()))
		if h.GetBattery() < 0{
			shortage := -(h.GetBattery()) 
			if p.GetBattery() < shortage{
				h.AddBlackout()
				h.ResetBattery()
			} else {
				p.WithdrawEnergy(-h.GetBattery())
			}

		}
		h.AddBattery(10*h.GetGG()) // change int to change capacity experiment_id_1 unlimited battery per house, no pool, no exchange
		donateTreshold := h.Getlast3DaysConsumption() // E3 
		if h.GetBattery() > donateTreshold{ // E3
			extraBattery := h.GetBattery() - donateTreshold // E3
			p.ContributeEnergy(extraBattery)
			h.AddBattery(-extraBattery)

		}
		PoolBattery = append(PoolBattery, p.GetBattery())
		

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
	pool := pool.NewPool(0)
	
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

	iterateDates(start_date, end_date, houses, pool)

	
	
	// used just for testing energy data stored in a house
	// for _, h := range houses {
	// 	if h.GetCustomer() == "11" {
	// 		err := SaveEnergyDataToCSV(h, "consumer_11_energydata.csv")
	// 		if err != nil {
	// 			fmt.Println("Error saving data:", err)
	// 		}
	// 		break
	// 	}
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

	
	
		
}



