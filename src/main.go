package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	_ "github.com/joho/godotenv/autoload"
	"googlemaps.github.io/maps"
)

func promptForAddress() string {
	fmt.Print("Enter Address: ")
	scan := bufio.NewScanner(os.Stdin)
	var addr string
	if (scan.Scan()) {
		addr = scan.Text()
	}
	return addr
}

func geocodeAddr(addr string) []maps.GeocodingResult {
	mapsKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	gm, err := maps.NewClient(maps.WithAPIKey(mapsKey))
	if err != nil {
		fmt.Println(err.Error())
		return make([]maps.GeocodingResult, 0)
	}
	r := &maps.GeocodingRequest{
		Address: addr,
	}
	res, err := gm.Geocode(context.Background(), r)
	if err != nil {
		fmt.Println(err.Error())
		return make([]maps.GeocodingResult, 0)
	}
	return res;
}

func deleteAddrAndSave() {

}


func main() {
// 	if (strings.ToLower(os.Args[1]) == "location") {
// fmt.Println(
// `Update Locations:
// [A]dd new location
// [R]emove location`)
// 		var in string
// 		fmt.Print(">")
// 		fmt.Scanln(&in)
// 		switch in {
// 		case "a":
// 		case "A":
// 			promptForAddrAndSave()

// 		case "r":
// 		case "R":
// 			deleteAddrAndSave()
// 		}

	
		
// 	} else {

// 	}

	addr := "12171 Beach Blvd, Jacksonville FL"
	loc := geocodeAddr(addr)
	if (len(loc) > 0) {

		// Get Point
		point, err := GetNWSPoint(loc[0].Geometry.Location.Lat, loc[0].Geometry.Location.Lng)

		if (err != nil) {
			fmt.Println("Could not get NWS point for current location")
			os.Exit(1)
		}
		
		// Get Forecast
		forecast, err := GetNWSForecastFromPoint(point)

		if (err != nil) {
			fmt.Println("Could not get forecast for point")
			os.Exit(1)
		}

		// Get Gridpoint
		gridpoint, err := GetNWSWeatherDetailFromPoint(point)
		if (err != nil) {
			fmt.Println("Could not get raw forecast data for point")
			os.Exit(0)
		}
		PrintForecast(forecast, point, gridpoint)
		
	}
}

func PrintForecast(f NWSForecast, l NWSPoint, g NWSWeatherDetail) {
	currentPeriod := f.Properties.Periods[0]
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgCyan)
	yellow.Printf("%s, %s\n", l.Properties.RelativeLocation.Properties.City, l.Properties.RelativeLocation.Properties.State)
	yellow.Printf("%s / %s\n", currentPeriod.Name, currentPeriod.StartTime.Format("01-02-2006 3:04PM"))
	yellow.Printf("%d°%s - %d%% chance of rain\n", currentPeriod.Temperature, currentPeriod.TemperatureUnit, currentPeriod.ProbabilityOfPrecipitation.Value)
	fmt.Println()
	blue.Println(currentPeriod.DetailedForecast)
	if (len(os.Args) > 1 && os.Args[1] == "-a") {
		fmt.Println("+---------------+------------------------------------+")
		fmt.Printf("| Temperature   | Cur %5s : Hi %6s  : Lo %6s |\n",
					g.Properties.Temperature.LastTempStr(),
					g.Properties.MaxTemperature.LastTempStr(),
					g.Properties.MinTemperature.LastTempStr())
		fmt.Printf("| Humidity/Rain | Hum %3d%%  : HIdx %4s  : COR %3d%%  |\n",
					g.Properties.RelativeHumidity.Last(),
					g.Properties.HeatIndex.LastTempStr(),
					g.Properties.ChanceOfRain.Last())
		fmt.Println("+---------------+------------------------------------+")

	}

	fmt.Println()
}
