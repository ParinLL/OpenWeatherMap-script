package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	geoURL      = "http://api.openweathermap.org/geo/1.0/direct"
	weatherURL  = "https://api.openweathermap.org/data/2.5/weather"
	forecastURL = "https://api.openweathermap.org/data/2.5/forecast"
)

// --- API response structs ---

type GeoResult struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type WeatherMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Humidity  int     `json:"humidity"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
}

type WeatherDesc struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}

type WeatherResponse struct {
	Name    string        `json:"name"`
	Main    WeatherMain   `json:"main"`
	Weather []WeatherDesc `json:"weather"`
	Wind    Wind          `json:"wind"`
}

type ForecastItem struct {
	Dt      int64         `json:"dt"`
	Main    WeatherMain   `json:"main"`
	Weather []WeatherDesc `json:"weather"`
	Wind    Wind          `json:"wind"`
	DtTxt   string        `json:"dt_txt"`
}

type ForecastResponse struct {
	List []ForecastItem `json:"list"`
	City struct {
		Name string `json:"name"`
	} `json:"city"`
}

func main() {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		fatal("OPENWEATHER_API_KEY env is required")
	}

	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "geo":
		if len(os.Args) < 3 {
			fatal("usage: owget geo <City>[,<Country>]")
		}
		cmdGeo(apiKey, os.Args[2])
	case "city":
		if len(os.Args) < 3 {
			fatal("usage: owget city <City>[,<Country>] [forecast]")
		}
		cmdCity(apiKey, os.Args[2], os.Args[3:])
	case "weather":
		lat, lon := parseLatLon(os.Args[2:])
		cmdWeather(apiKey, lat, lon)
	case "forecast":
		lat, lon := parseLatLon(os.Args[2:])
		cmdForecast(apiKey, lat, lon)
	default:
		// default: treat args as lat lon for weather
		lat, lon := parseLatLon(os.Args[1:])
		cmdWeather(apiKey, lat, lon)
	}
}

func parseLatLon(args []string) (float64, float64) {
	if len(args) < 2 {
		fatal("usage: owget [weather|forecast] <lat> <lon>")
	}
	lat, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fatal("invalid lat: " + args[0])
	}
	lon, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fatal("invalid lon: " + args[1])
	}
	return lat, lon
}

func cmdGeo(apiKey, query string) {
	url := fmt.Sprintf("%s?q=%s&limit=5&appid=%s", geoURL, url.QueryEscape(query), apiKey)
	body := httpGet(url)

	var results []GeoResult
	mustUnmarshal(body, &results)

	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}
	for _, r := range results {
		fmt.Printf("📍 %s, %s  (lat=%.4f, lon=%.4f)\n", r.Name, r.Country, r.Lat, r.Lon)
	}
}

func cmdCity(apiKey, query string, extra []string) {
	reqURL := fmt.Sprintf("%s?q=%s&limit=1&appid=%s", geoURL, url.QueryEscape(query), apiKey)
	body := httpGet(reqURL)

	var results []GeoResult
	mustUnmarshal(body, &results)

	if len(results) == 0 {
		fatal("city not found: " + query)
	}

	r := results[0]
	fmt.Printf("📍 %s, %s  (lat=%.4f, lon=%.4f)\n\n", r.Name, r.Country, r.Lat, r.Lon)

	if len(extra) > 0 && extra[0] == "forecast" {
		cmdForecast(apiKey, r.Lat, r.Lon)
	} else {
		cmdWeather(apiKey, r.Lat, r.Lon)
	}
}


func cmdWeather(apiKey string, lat, lon float64) {
	url := fmt.Sprintf("%s?lat=%.4f&lon=%.4f&units=metric&appid=%s", weatherURL, lat, lon, apiKey)
	body := httpGet(url)

	var w WeatherResponse
	mustUnmarshal(body, &w)

	desc := ""
	if len(w.Weather) > 0 {
		desc = w.Weather[0].Description
	}

	fmt.Printf("🌡️  %s — %s\n", w.Name, desc)
	fmt.Printf("   Temp: %.1f°C (feels %.1f°C)  Low %.1f / High %.1f\n",
		w.Main.Temp, w.Main.FeelsLike, w.Main.TempMin, w.Main.TempMax)
	fmt.Printf("   Humidity: %d%%  Wind: %.1f m/s\n", w.Main.Humidity, w.Wind.Speed)
}

func cmdForecast(apiKey string, lat, lon float64) {
	url := fmt.Sprintf("%s?lat=%.4f&lon=%.4f&units=metric&appid=%s", forecastURL, lat, lon, apiKey)
	body := httpGet(url)

	var f ForecastResponse
	mustUnmarshal(body, &f)

	fmt.Printf("📅 5-Day Forecast — %s\n", f.City.Name)
	fmt.Println(strings.Repeat("─", 60))

	lastDate := ""
	for _, item := range f.List {
		t := time.Unix(item.Dt, 0).In(time.FixedZone("CST", 8*3600))
		date := t.Format("01/02 (Mon)")
		hour := t.Format("15:04")

		if date != lastDate {
			fmt.Printf("\n  %s\n", date)
			lastDate = date
		}

		desc := ""
		if len(item.Weather) > 0 {
			desc = item.Weather[0].Description
		}
		fmt.Printf("    %s  %5.1f°C  💧%d%%  💨%.1fm/s  %s\n",
			hour, item.Main.Temp, item.Main.Humidity, item.Wind.Speed, desc)
	}
}

func httpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fatal("http error: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fatal("read error: " + err.Error())
	}
	if resp.StatusCode != 200 {
		fatal(fmt.Sprintf("API %d: %s", resp.StatusCode, string(body)))
	}
	return body
}

func mustUnmarshal(data []byte, v any) {
	if err := json.Unmarshal(data, v); err != nil {
		fatal("json error: " + err.Error())
	}
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

Env:
  OPENWEATHER_API_KEY                Required

Examples:
  owget 24.9575 121.5105
  owget forecast 25.0287 121.5052
  owget city Taipei,TW
  owget city Taipei,TW forecast
  owget geo Ankang,TW
`)
	os.Exit(1)
}

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}
