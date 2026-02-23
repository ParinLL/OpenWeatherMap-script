package cmd

import (
	"fmt"
	"time"

	"owget/api"
)

func Weather(apiKey string, lat, lon float64, detail bool) {
	reqURL := fmt.Sprintf("%s?lat=%.4f&lon=%.4f&units=metric&appid=%s", api.WeatherURL, lat, lon, apiKey)
	body := api.HTTPGet(reqURL)

	var w api.WeatherResponse
	api.Unmarshal(body, &w)

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
