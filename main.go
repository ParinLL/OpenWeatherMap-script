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

var debug bool
var detail bool

// --- API response structs ---

type GeoResult struct {
	Name    string            `json:"name"`
	Country string            `json:"country"`
	State   string            `json:"state"`
	Lat     float64           `json:"lat"`
	Lon     float64           `json:"lon"`
	Local   map[string]string `json:"local_names"`
}

type WeatherMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Humidity  int     `json:"humidity"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type WeatherDesc struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

type Rain struct {
	OneH   float64 `json:"1h"`
	ThreeH float64 `json:"3h"`
}

type WeatherResponse struct {
	Name       string        `json:"name"`
	Main       WeatherMain   `json:"main"`
	Weather    []WeatherDesc `json:"weather"`
	Wind       Wind          `json:"wind"`
	Clouds     Clouds        `json:"clouds"`
	Sys        Sys           `json:"sys"`
	Rain       *Rain         `json:"rain"`
	Visibility int           `json:"visibility"`
	Dt         int64         `json:"dt"`
	Timezone   int           `json:"timezone"`
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

	// check for flags
	args := []string{}
	for _, a := range os.Args[1:] {
		switch a {
		case "--debug":
			debug = true
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
			fatal("usage: owget geo <City>[,<Country>]")
		}
		cmdGeo(apiKey, args[1])
	case "city":
		if len(args) < 2 {
			fatal("usage: owget city <City>[,<Country>] [forecast]")
		}
		cmdCity(apiKey, args[1], args[2:])
	case "weather":
		lat, lon := parseLatLon(args[1:])
		cmdWeather(apiKey, lat, lon)
	case "forecast":
		lat, lon := parseLatLon(args[1:])
		cmdForecast(apiKey, lat, lon)
	default:
		lat, lon := parseLatLon(args)
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

// encodeQuery URL-encodes each comma-separated part of the query
// while preserving commas as-is (OpenWeatherMap uses commas as delimiters).
func encodeQuery(query string) string {
	parts := strings.Split(query, ",")
	for i, p := range parts {
		parts[i] = url.QueryEscape(strings.TrimSpace(p))
	}
	return strings.Join(parts, ",")
}

func cmdGeo(apiKey, query string) {
	reqURL := fmt.Sprintf("%s?q=%s&limit=5&appid=%s", geoURL, encodeQuery(query), apiKey)
	body := httpGet(reqURL)

	var results []GeoResult
	mustUnmarshal(body, &results)

	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}
	for _, r := range results {
		fmt.Printf("📍 %s, %s  (lat=%.4f, lon=%.4f)\n", r.Name, r.Country, r.Lat, r.Lon)
		if detail {
			if r.State != "" {
				fmt.Printf("   State: %s\n", r.State)
			}
			if len(r.Local) > 0 {
				names := []string{}
				for lang, name := range r.Local {
					if lang == "ascii" || lang == "feature_name" {
						continue
					}
					names = append(names, fmt.Sprintf("%s:%s", lang, name))
				}
				fmt.Printf("   Local: %s\n", strings.Join(names, ", "))
			}
			fmt.Println()
		}
	}
}

func cmdCity(apiKey, query string, extra []string) {
	reqURL := fmt.Sprintf("%s?q=%s&limit=1&appid=%s", geoURL, encodeQuery(query), apiKey)
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
	reqURL := fmt.Sprintf("%s?lat=%.4f&lon=%.4f&units=metric&appid=%s", weatherURL, lat, lon, apiKey)
	body := httpGet(reqURL)

	var w WeatherResponse
	mustUnmarshal(body, &w)

	desc := ""
	if len(w.Weather) > 0 {
		desc = w.Weather[0].Description
	}

	tz := time.FixedZone("local", w.Timezone)

	fmt.Printf("🌡️  %s — %s\n", w.Name, desc)
	fmt.Printf("   Temp: %.1f°C (feels %.1f°C)  Low %.1f / High %.1f\n",
		w.Main.Temp, w.Main.FeelsLike, w.Main.TempMin, w.Main.TempMax)
	fmt.Printf("   Humidity: %d%%  Wind: %.1f m/s\n", w.Main.Humidity, w.Wind.Speed)

	if detail {
		fmt.Printf("   Pressure: %d hPa", w.Main.Pressure)
		if w.Main.SeaLevel > 0 {
			fmt.Printf("  Sea: %d hPa  Grnd: %d hPa", w.Main.SeaLevel, w.Main.GrndLevel)
		}
		fmt.Println()
		fmt.Printf("   Wind: %.1f m/s  Dir: %d°", w.Wind.Speed, w.Wind.Deg)
		if w.Wind.Gust > 0 {
			fmt.Printf("  Gust: %.1f m/s", w.Wind.Gust)
		}
		fmt.Println()
		fmt.Printf("   Clouds: %d%%  Visibility: %dm\n", w.Clouds.All, w.Visibility)
		if w.Rain != nil {
			if w.Rain.OneH > 0 {
				fmt.Printf("   Rain (1h): %.1fmm\n", w.Rain.OneH)
			}
			if w.Rain.ThreeH > 0 {
				fmt.Printf("   Rain (3h): %.1fmm\n", w.Rain.ThreeH)
			}
		}
		if w.Sys.Sunrise > 0 {
			sunrise := time.Unix(w.Sys.Sunrise, 0).In(tz).Format("15:04")
			sunset := time.Unix(w.Sys.Sunset, 0).In(tz).Format("15:04")
			fmt.Printf("   Sunrise: %s  Sunset: %s\n", sunrise, sunset)
		}
		observed := time.Unix(w.Dt, 0).In(tz).Format("2006-01-02 15:04")
		fmt.Printf("   Observed: %s (UTC%+d)\n", observed, w.Timezone/3600)
	}
}

func cmdForecast(apiKey string, lat, lon float64) {
	reqURL := fmt.Sprintf("%s?lat=%.4f&lon=%.4f&units=metric&appid=%s", forecastURL, lat, lon, apiKey)
	body := httpGet(reqURL)

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
		if detail {
			fmt.Printf("           feels %.1f°C  low %.1f / high %.1f  pressure %dhPa  wind dir %d°\n",
				item.Main.FeelsLike, item.Main.TempMin, item.Main.TempMax, item.Main.Pressure, item.Wind.Deg)
		}
	}
}

func httpGet(reqURL string) []byte {
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] GET %s\n", reqURL)
	}
	resp, err := http.Get(reqURL)
	if err != nil {
		fatal("http error: " + err.Error())
	}
	defer resp.Body.Close()
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] HTTP %d\n", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fatal("read error: " + err.Error())
	}
	if debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Body: %s\n", string(body))
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

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}
