package gosurf

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestCalculate(t *testing.T) {
	type args struct {
		spot    Spot
		spotMap map[string]map[int]Spot
	}

	testCases := []struct {
		desc   string
		args   args
		result map[string]Spot
	}{
		{
			desc: "should calculate tide position all successfully",
			args: args{
				spot: Spot{
					Name: "test",
					Lng:  0.1,
					Lat:  0.1,
					Tide: Tide{
						Position: "all",
					},
					Period: 10,
				},
				spotMap: map[string]map[int]Spot{
					"2001-01-01": {
						0: {
							Period: 11,
						},
					},
					"2001-01-02": {
						1: {
							Period: 0,
						},
					},
				},
			},
			result: map[string]Spot{
				"2001-01-01": {
					Period: 11,
				},
			},
		},
		{
			desc: "should calculate tide position low successfully",
			args: args{
				spot: Spot{
					Name: "test",
					Lng:  0.1,
					Lat:  0.1,
					Tide: Tide{
						Position: "low",
					},
					Period: 10,
				},
				spotMap: map[string]map[int]Spot{
					"2001-01-01": {
						time.Now().Hour() - 1: {
							Period: 11,
							Tide: Tide{
								Position:   "low",
								TimeBefore: 2,
								TimeAfter:  2,
								Low:        time.Now(),
							},
						},
					},
					"2001-01-02": {
						time.Now().Hour() - 3: {
							Period: 11,
							Tide: Tide{
								Position:   "low",
								TimeBefore: 2,
								TimeAfter:  2,
								Low:        time.Now(),
							},
						},
					},
				},
			},
			result: map[string]Spot{
				"2001-01-01": {
					Period: 11,
					Tide: Tide{
						Position:   "low",
						TimeBefore: 2,
						TimeAfter:  2,
						Low:        time.Now(),
					},
				},
			},
		},
		{
			desc: "should calculate tide position high successfully",
			args: args{
				spot: Spot{
					Name: "test",
					Lng:  0.1,
					Lat:  0.1,
					Tide: Tide{
						Position: "high",
					},
					Period: 10,
				},
				spotMap: map[string]map[int]Spot{
					"2001-01-01": {
						time.Now().Hour() - 1: {
							Period: 11,
							Tide: Tide{
								Position:   "high",
								TimeBefore: 2,
								TimeAfter:  2,
								High:       time.Now(),
							},
						},
					},
					"2001-01-02": {
						time.Now().Hour() - 3: {
							Period: 11,
							Tide: Tide{
								Position:   "high",
								TimeBefore: 2,
								TimeAfter:  2,
								High:       time.Now(),
							},
						},
					},
				},
			},
			result: map[string]Spot{
				"2001-01-01": {
					Period: 11,
					Tide: Tide{
						Position:   "high",
						TimeBefore: 2,
						TimeAfter:  2,
						High:       time.Now(),
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			result := calculate(tt.args.spot, tt.args.spotMap)
			assert.Equal(t, tt.result, result)
		})
	}
}
