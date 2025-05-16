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
