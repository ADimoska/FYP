package main

import (
	"CES/house"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// Generates Histogram for total Blackouts per house
func CallPythonToGenerateHistogram(houses []*house.House) error {
	// totalBlackouts := 0
	// var blackouts []int
	// zeroBlackoutCount := 0

	// for _, h := range houses {
	// 	count := h.GetBlackouts()
	// 	blackouts = append(blackouts, count)
	// 	if count == 0 {
	// 		zeroBlackoutCount++
	// 	}
	// 	totalBlackouts += count

	// }
	// fmt.Println("Total number of blackouts, python func result:", totalBlackouts)
	// fmt.Println("Number of houses with 0 blackouts:", zeroBlackoutCount)
	// totalBlackouts_global = totalBlackouts
	// zeroBlackoutCount_global = zeroBlackoutCount
	blackouts := CountTotalBlackouts(houses)

	
	data, err := json.Marshal(blackouts)
	if err != nil {
		return err
	}

	cmd := exec.Command("python3", "graph.py") 
	cmd.Stdin = bytes.NewReader(data)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("python error: %v, stderr: %s", err, stderr.String())
	}

	fmt.Println("Python plot saved as blackout_distribution.png")
	return nil
}

func CallPythonToPlotPoolBattery(PoolBattery []float64) error {
	// Serialize the PoolBattery slice into JSON
	data, err := json.Marshal(PoolBattery)
	if err != nil {
		return fmt.Errorf("failed to marshal PoolBattery: %v", err)
	}

	// Prepare the command to run the Python script
	cmd := exec.Command("python3", "plot_pool_battery.py")
	cmd.Stdin = bytes.NewReader(data)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the Python script
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("python error: %v, stderr: %s", err, stderr.String())
	}

	fmt.Println("Python plot saved as pool_battery_plot.png")
	return nil
}


func CountTotalBlackouts(houses []*house.House) []int{
	totalBlackouts := 0
	var blackouts []int
	zeroBlackoutCount := 0

	for _, h := range houses {
		count := h.GetBlackouts()
		blackouts = append(blackouts, count)
		if count == 0 {
			zeroBlackoutCount++
		}
		totalBlackouts += count

	}
	fmt.Println("Total number of blackouts, python func result:", totalBlackouts)
	fmt.Println("Number of houses with 0 blackouts:", zeroBlackoutCount)
	totalBlackouts_global = totalBlackouts
	zeroBlackoutCount_global = zeroBlackoutCount
	return blackouts
}

func CountBlackoutsByCommunity(houses []*house.House) (int, int) {
	community1Blackouts := 0
	community2Blackouts := 0

	zeroBlackoutCountC1 := 0
	zeroBlackoutCountC2 := 0


	for _, h := range houses {
		count := h.GetBlackouts()
		community := h.GetCommunityID()

		if community == 1 {
			if count == 0 {
				zeroBlackoutCountC1++
			}
			community1Blackouts += count
		} else if community == 2 {
			if count == 0 {
				zeroBlackoutCountC2++
			}
			community2Blackouts += count
		}
	}

	fmt.Println("Total blackouts in Community 1:", community1Blackouts)
	fmt.Println("Total blackouts in Community 2:", community2Blackouts)
	totalBlackoutsC1_global = community1Blackouts
	totalBlackoutsC2_global = community2Blackouts
	zeroBlackoutCountC1_global = zeroBlackoutCountC1
	zeroBlackoutCountC2_global = zeroBlackoutCountC2

	return community1Blackouts, community2Blackouts
}

