package main

import (
	"fmt"
	"os"
	"strconv"

	"owget/api"
	"owget/cmd"
)

func main() {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		api.Fatal("OPENWEATHER_API_KEY env is required")
	}

	if len(os.Args) < 2 {
		usage()
	}

	// parse flags
	var detail bool
	args := []string{}
	for _, a := range os.Args[1:] {
		switch a {
		case "--debug":
			api.Debug = true
		case "--detail":
			detail = true
		default:
			args = append(args, a)
		}
	}

	if len(args) == 0 {
		usage()
	}

	switch args[0] {
	case "geo":
		if len(args) < 2 {
			api.Fatal("usage: owget geo <City>[,<Country>]")
		}
		cmd.Geo(apiKey, args[1], detail)
	case "city":
		if len(args) < 2 {
			api.Fatal("usage: owget city <City>[,<Country>] [forecast]")
		}
		cmd.City(apiKey, args[1], args[2:], detail)
	case "weather":
		lat, lon := parseLatLon(args[1:])
		cmd.Weather(apiKey, lat, lon, detail)
	case "forecast":
		lat, lon := parseLatLon(args[1:])
		cmd.Forecast(apiKey, lat, lon, detail)
	default:
		lat, lon := parseLatLon(args)
		cmd.Weather(apiKey, lat, lon, detail)
	}
}

func parseLatLon(args []string) (float64, float64) {
	if len(args) < 2 {
		api.Fatal("usage: owget [weather|forecast] <lat> <lon>")
	}
	lat, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		api.Fatal("invalid lat: " + args[0])
	}
	lon, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		api.Fatal("invalid lon: " + args[1])
	}
	return lat, lon
}

func usage() {
	fmt.Fprintf(os.Stderr, `owget — OpenWeatherMap CLI

Usage:
  owget <lat> <lon>                  Current weather (shortcut)
  owget weather <lat> <lon>          Current weather
  owget forecast <lat> <lon>         5-day / 3-hour forecast
  owget city <City>[,<Country>]      Current weather by city name
  owget city <City>[,<Country>] forecast
                                     5-day forecast by city name
  owget geo <City>[,<Country>]       Search location

Flags:
  --detail                           Show detailed weather information
  --debug                            Show HTTP request/response details

Env:
  OPENWEATHER_API_KEY                Required

Examples:
  owget 24.9575 121.5105
  owget forecast 25.0287 121.5052
  owget city Taipei,TW
  owget city "New York,NY,US"
  owget geo "New York,US" --debug
`)
	os.Exit(1)
}
