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
var averageThresholds []float64 //used for autoE2E3
var totalBlackouts_global_array []int
var totalBlackoutsC1_global_array []int
var totalBlackoutsC2_global_array []int
var zeroBlackoutCount_global_array []int
var zeroBlackoutCountC1_global_array []int
var zeroBlackoutCountC2_global_array []int
var totalBlackouts_global int = 0 // used for autoE2E3
var totalBlackoutsC1_global int = 0 
var totalBlackoutsC2_global int = 0 
var zeroBlackoutCount_global int = 0 // used for autoE2E3
var zeroBlackoutCountC1_global int = 0
var zeroBlackoutCountC2_global int = 0
var PoolBattery1 []float64
// iterateDates takes two date strings (start and end) in the format "2/01/2006",
// parses them into time.Time values, and prints each date from the start to the end, inclusive.
// It validates that both dates are correctly formatted and that the start date is not after the end date.
// If any validation fails, it prints an error message and exits the function.
func iterateDates(startDateStr, endDateStr string, houses []*house.House, p1, p2 *pool.Pool, input int) { //E2E3auto
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
		executeDate(date_str, houses, p1, p2, input) //E2E3auto
	}
	//GenerateBlackoutChart(houses, "blackouts.html")
	
	//CallPythonToGenerateHistogram(houses)
	CountTotalBlackouts(houses)
	CountBlackoutsByCommunity(houses)
	CallPythonToPlotPoolBattery(PoolBattery1)
	fmt.Println("Average threshold:",(total3DaySum/total3DayCount)) //E3 
	
}

var total3DaySum float64 = 0.0 //E3
var total3DayCount float64 = 0 //E3
func executeDate (date string, houses []*house.House, p1, p2 *pool.Pool, input int) {
	timeArray := []string{"0:30","1:00","1:30","2:00","2:30","3:00","3:30","4:00","4:30","5:00","5:30","6:00","6:30","7:00","7:30","8:00","8:30","9:00","9:30","10:00","10:30","11:00","11:30","12:00","12:30","13:00","13:30","14:00","14:30","15:00","15:30","16:00","16:30","17:00","17:30","18:00","18:30","19:00","19:30","20:00","20:30","21:00","21:30","22:00","22:30","23:00","23:30","0:00"}	
	// E3 start =================================================================
	// if date == "1/07/2012" || date == "2/07/2012" || date == "3/07/2012"{
	// 	for _, h := range houses {
	// 	h.SetNext3Days(date, float64(input))
	// 	total3DaySum += h.Getlast3DaysConsumption()
	// 	total3DayCount++
	// 	}
	// } else {
	// 	for _, h := range houses {
	// 		h.SetLast3Days(date, float64(input))
	// 		total3DaySum += h.Getlast3DaysConsumption()
	// 		total3DayCount++
	// 		}
	// }
	// E3 end =================================================================
	
	for _, time := range timeArray {
		executeTime(date, time, houses, p1, p2, input)

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

func executeTime(date, time string, houses []*house.House, p1, p2 *pool.Pool, input int) {
	for _, h := range houses {
		var p *pool.Pool
		// var po *pool.Pool		//E4
		c := h.GetCommunityID()
		if c == 1 {
			p = p1
		} else {
			p = p2
		}
		// if p == p1{
		// 	po = p2
		// } else {
		// 	po = p1
		// }
		h.GetCurrentEnergy(date, time)
		h.AddBattery(-(h.GetGC()))
		h.AddBattery(-(h.GetCL()))
		if h.GetBattery() < 0{
			shortage := -(h.GetBattery()) 
			if p.GetBattery() < shortage{
				// if po.GetBattery() < shortage{
					h.AddBlackout()
					h.ResetBattery()
				// } else {
					// po.GiveToOtherPool(-h.GetBattery(), p)
					// p.GiveEnergy(-h.GetBattery())
				// }
				
			} else {
				p.GiveEnergy(-h.GetBattery())
			}

		}
		multiGG := getGGMultiplier(c, time) //E4
		h.AddBattery(multiGG*h.GetGG()) // change int to change capacity experiment_id_1 unlimited battery per house, no pool, no exchange
		donateTreshold := 56.00 // E3 E4
		// donateTreshold := input
		if h.GetBattery() > donateTreshold{ // E3
			extraBattery := h.GetBattery() - donateTreshold // E3
			p.AcceptEnergy(extraBattery)
			if p.GetBattery() > p.GetCapacity(){
				p.SetEnergy(p.GetCapacity())

				// poolExtraBattery := p.GetBattery() - p.GetCapacity()
				// p.GiveToOtherPool(poolExtraBattery, po)
			}
			h.AddBattery(-extraBattery)

		}
		PoolBattery1 = append(PoolBattery1, p2.GetBattery())
		

	}
}

func runSimulation(input int) {

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
	pool1 := pool.NewPool(0, 56*2*49)
	pool2 := pool.NewPool(0, 56*2*65)
	
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

	// houses = reverseHouses(houses)

	for _, h := range houses{
		fmt.Printf("House %s:%s\n", h.GetCustomer(), h.GetCity())
	}

	iterateDates(start_date, end_date, houses, pool1, pool2, input)

	fmt.Print(countHousesByCommunity(houses))


	
	
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
	// fmt.Println("--")
    // fmt.Println(dataStore["Sydney"][2012]["Jan"][1])
	// fmt.Println(dataStore["Sydney"][2013]["Jan"][1])
	// fmt.Println(dataStore["Gosford"][2012]["Jan"][1])
	// fmt.Println(dataStore["Gosford"][2013]["Jan"][1])
	// fmt.Println(dataStore["Newcastle"][2012]["Jan"][1])
	// fmt.Println(dataStore["Newcastle"][2013]["Jan"][1])
	// fmt.Println(dataStore["Cessnock"][2012]["Jan"][1])
	// fmt.Println(dataStore["Cessnock"][2013]["Jan"][1])
	


	// Store the average threshold
	if total3DayCount > 0 {
		average := total3DaySum / total3DayCount
		averageThresholds = append(averageThresholds, average)
		fmt.Println("Average threshold:", average)
	} else {
		fmt.Println("No valid data for threshold calculation.")
	}


	totalBlackouts_global_array = append(totalBlackouts_global_array, totalBlackouts_global)
	totalBlackoutsC1_global_array = append(totalBlackoutsC1_global_array, totalBlackoutsC1_global)
	totalBlackoutsC2_global_array = append(totalBlackoutsC2_global_array, totalBlackoutsC2_global)
	zeroBlackoutCount_global_array = append(zeroBlackoutCount_global_array, zeroBlackoutCount_global)
	zeroBlackoutCountC1_global_array = append(zeroBlackoutCountC1_global_array, zeroBlackoutCountC1_global)
	zeroBlackoutCountC2_global_array = append(zeroBlackoutCountC2_global_array, zeroBlackoutCountC2_global)
	// Reset global counters before each run
	total3DaySum = 0.0
	total3DayCount = 0.0
	totalBlackouts_global = 0
	totalBlackoutsC1_global = 0
	totalBlackoutsC2_global = 0
	zeroBlackoutCount_global  = 0 
	zeroBlackoutCountC1_global  = 0 
	zeroBlackoutCountC2_global  = 0 
	
	// cityCounts := countHousesByCity(houses)
	// fmt.Println("House counts by city:")
	// for city, count := range cityCounts {
	// 	fmt.Printf("%s: %d houses\n", city, count)
	// }
	// c1, c2 := countHousesByCommunity(houses)
	// fmt.Printf("Community 1: %d houses\n", c1)
	// fmt.Printf("Community 2: %d houses\n", c2)
		
}




func main() {
	for i := 0; i < 1; i++ { 
		fmt.Printf("Running simulation #%d\n", i+1)
		runSimulation(i)
	}

	// average_thresholds := []float64{0.00, 56.03, 112.06, 168.09, 224.12, 280.15, 336.18, 392.21, 448.24, 504.27, 560.30, 616.33, 672.36, 728.39, 784.42, 840.44, 896.47, 952.50, 1008.53, 1064.56, 1120.59, 1176.62, 1232.65, 1288.68, 1344.71, 1400.74, 1456.77, 1512.80, 1568.83, 1624.86, 1680.89, 1736.92, 1792.95, 1848.98, 1905.01, 1961.04, 2017.07, 2073.10, 2129.13, 2185.16, 2241.19}
	// average_thresholds := []float64{0.00}
	// for i, avg := range average_thresholds {
	// 	fmt.Printf("Running simulation #%d\n", i+1)
	// 	runSimulation(avg)
	// }

	

	fmt.Println("All average thresholds from simulations:")
	for i, avg := range averageThresholds {
		fmt.Printf("Run #%d: %.2f, %.2d, %.2d\n", i+1, avg, totalBlackouts_global_array[i], zeroBlackoutCount_global_array[i] )
	}

	fmt.Print("Average Thresholds: ")
	for i, v := range averageThresholds {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%.2f", v)
	}
	fmt.Println()

	fmt.Print("Total Blackouts: ")
	for i, v := range totalBlackouts_global_array {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println()


	fmt.Print("Total Blackouts C1: ")
	for i, v := range totalBlackoutsC1_global_array {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println()


	fmt.Print("Total Blackouts C2: ")
	for i, v := range totalBlackoutsC2_global_array {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println()


	fmt.Print("Zero Blackout Houses: ")
	for i, v := range zeroBlackoutCount_global_array {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println()

	fmt.Print("Zero Blackout Houses C1: ")
	for i, v := range zeroBlackoutCountC1_global_array {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println()

	fmt.Print("Zero Blackout Houses C2: ")
	for i, v := range zeroBlackoutCountC2_global_array {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%d", v)
	}
	fmt.Println()
}
