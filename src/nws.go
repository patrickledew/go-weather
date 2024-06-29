package main

import (
	"fmt"
	"strings"
	"time"
)

/** Helpful methods **/
// Get NWS's point object, for a location in their grid
func GetNWSPoint(lat float64, lng float64) (NWSPoint, error) {
	var point NWSPoint
	_, err := get_api_into(&point, "https://api.weather.gov/points/%f,%f", lat, lng)
	return point, err
}

// Get a forecast for a specific point
func GetNWSForecastFromPoint(point NWSPoint) (NWSForecast, error) {
	var forecast NWSForecast
	_, err := get_api_into(&forecast, point.Properties.Forecast)
	return forecast, err
}

func GetNWSWeatherDetailFromPoint(point NWSPoint) (NWSWeatherDetail, error) {
	var gridpoint NWSWeatherDetail
	_, err := get_api_into(&gridpoint, "https://api.weather.gov/gridpoints/%s/%d,%d", point.Properties.GridId, point.Properties.GridX, point.Properties.GridX)
	return gridpoint, err
}


/** API MODELS **/

/* Full NWS API docs https://www.weather.gov/documentation/services-web-api */


// Represents a latitude/longitude point
type NWSPoint struct {
	Properties struct 
	{
		GridX int
		GridY int
		Forecast string
		ForecastHourly string
		RelativeLocation struct {
			Properties struct {
				City string
				State string
			}
		}
		GridId string
	}
}

// Contains a basic forecast for an area.
type NWSForecast struct {
	Properties struct {
		Periods []struct {
			Name string
			StartTime time.Time
			Temperature int
			TemperatureUnit string
			ProbabilityOfPrecipitation struct {
				Value int
			}
			DetailedForecast string
		}
	}
}

// Contains raw numerical forecast data across several different metrics.
type NWSWeatherDetail struct {
	Properties struct {
		UpdateTime time.Time
		Temperature NWSTempSeries
		DewPoint NWSDataSeries[float64]
		MaxTemperature NWSTempSeries
		MinTemperature NWSTempSeries
		RelativeHumidity NWSDataSeries[int]
		ApparentTemperature NWSTempSeries
		HeatIndex NWSTempSeries
		WindChill NWSDataSeries[float64]
		SkyCover NWSDataSeries[int]
		WindDirection NWSDataSeries[int]
		WindSpeed NWSDataSeries[float64]
		WindGust NWSDataSeries[float64]
		ChanceOfRain NWSDataSeries[int] `json:"probabilityOfPrecipitation"`
		Weather NWSDataSeries[[]NWSWeatherInfo]
	}
}

type NWSWeatherInfo struct {
	Coverage string // e.g slight_chance, chance, liekely
	Weather string // weather event, e.g. rain_showers, thunderstorms
	Intensity string // e.g moderate
}

type NWSDataPoint [T any] struct {
	
	/* The time and duration the data point covers. In ISO 8601 fmt,
	e.g. 2024-06-29T01:00:00+00:00/PT11H

	The duration at the end is expressed using the following format, where (n) is replaced by the value for each of the date and time elements that follow the (n):

    P(n)Y(n)M(n)DT(n)H(n)M(n)S
	*/
	ValidTime NWSValidTime
	Value T // Usually a number, but could be a more complex data structure for some series
}

type NWSDataSeries [T any] struct {
	/* 
	One of:

	wmoUnit:degC,
	wmoUnit:degF,
	wmoUnit:percent,
	wmoUnit:degree_(angle),
	wmoUnit:km_h-1
	wmoUnit:mm
	wmoUnit:m
	nwsUnit:s
	might be more idk
	*/ 
	Units string `json:"uom"`
	Values []NWSDataPoint[T]
}

type NWSValidTime string

func (timeAndDuration NWSValidTime) Time() time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05-0700", strings.Split(string(timeAndDuration), "/")[0])

	return t
}



func (series NWSDataSeries[T]) Last() T {
	return series.Values[0].Value
}

type NWSTempSeries NWSDataSeries[float64]

func (series NWSTempSeries) LastTemp() float64 {
	tempDisplay := series.Values[0].Value

	if (series.Units == "wmoUnit:degC" && TEMP_UNITS == "F") {
		tempDisplay *= 1.8
		tempDisplay += 32
	}
	if (series.Units == "wmoUnit:degF" && TEMP_UNITS == "C") {
		tempDisplay -= 32
		tempDisplay /= 1.8
	}
	return tempDisplay
}
func (series NWSTempSeries) LastTempStr() string {
	return fmt.Sprintf("%.f°%s", series.LastTemp(), TEMP_UNITS)
}

const TEMP_UNITS = "F"