package main

import (
	"CES/house"
	
	"fmt"	
	"os"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	//"github.com/go-echarts/go-echarts/v2/types"

	
)

// func GenerateBlackoutChart(houses []*house.House, outputFile string) error {
// 	bar := charts.NewBar()
// 	var houseIDs []string
// 	var blackoutCounts []opts.BarData
// 	var blackoutCounts_ []int

// 	for _, h := range houses {
// 		fmt.Println(h.GetCustomer(), ":", h.GetBlackouts())
// 		houseIDs = append(houseIDs, h.GetCustomer())
// 		blackoutCounts = append(blackoutCounts, opts.BarData{Value: h.GetBlackouts()})
// 		blackoutCounts_ = append(blackoutCounts_, h.GetBlackouts())
// 	}

// 	totalBlackouts := 0
// 	for _, count := range blackoutCounts_ {
// 		totalBlackouts += count
// 	}
// 	fmt.Println("Total blackout count:", totalBlackouts)

// 	// Set chart data
// 	bar.SetGlobalOptions(
// 		charts.WithTitleOpts(opts.Title{Title: "Blackouts per House"}),
// 	)
// 	bar.SetXAxis(houseIDs).AddSeries("Blackouts", blackoutCounts)

// 	// Render chart to HTML file
// 	f, err := os.Create(outputFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	err = bar.Render(f)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("Total houses:", len(houses))
// 	return nil
// }

func GenerateBlackoutChart(houses []*house.House, outputFile string) error {
	bar := charts.NewBar()
	var houseIDs []string
	var blackoutCounts []opts.BarData
	totalBlackouts := 0

	for _, h := range houses {
		customer := h.GetCustomer()
		blackouts := h.GetBlackouts()

		fmt.Println(customer, ":", blackouts)
		houseIDs = append(houseIDs, customer)
		blackoutCounts = append(blackoutCounts, opts.BarData{
			Value: blackouts,
			Label: &opts.Label{
				Show:     opts.Bool(true),
				Formatter: customer,
			},
		})
		totalBlackouts += blackouts
	}

	fmt.Println("Total blackout count:", totalBlackouts)

	// Set chart options
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Blackouts per House"}),
	)

	// Set data and enable bar labels
	bar.SetXAxis(houseIDs).AddSeries("Blackouts", blackoutCounts).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     opts.Bool(true),
				Position: "top",
			}),
			charts.WithBarChartOpts(opts.BarChart{
				// BarWidth: "40%", // ↓ Reduce width to increase space
				BarCategoryGap: "130%",
			}),
		)

	// Render chart to HTML file
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := bar.Render(f); err != nil {
		return err
	}

	fmt.Println("Total houses:", len(houses))
	return nil
}