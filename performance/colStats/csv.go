package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"
)

type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0

	for _, v := range data {
		sum += v
	}

	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func min(data []float64) float64 {
	min := math.Inf(1)

	for _, datum := range data {
		if datum < min {
			min = datum
		}
	}

	if min == math.Inf(1) {
		return 0
	}

	return min
}

func max(data []float64) float64 {
	max := math.Inf(-1)

	for _, datum := range data {
		if datum > max {
			max = datum
		}
	}

	if max == math.Inf(-1) {
		return 0
	}

	return max
}

func csv2float(r io.Reader, column int) ([]float64, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	// Adjust for 0 based index
	column--

	var data []float64

	for i := 0; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read data from file: %w", err)
		}

		// Discard the row with the header
		if i == 0 {
			continue
		}

		if len(row) <= column {
			return nil,
				fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		data = append(data, v)
	}

	return data, nil
}
