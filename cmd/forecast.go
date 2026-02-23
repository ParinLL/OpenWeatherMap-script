package cmd

import (
	"fmt"
	"strings"
	"time"

	"owget/api"
)

func Forecast(apiKey string, lat, lon float64, detail bool) {
	reqURL := fmt.Sprintf("%s?lat=%.4f&lon=%.4f&units=metric&appid=%s", api.ForecastURL, lat, lon, apiKey)
	body := api.HTTPGet(reqURL)

	var f api.ForecastResponse
	api.Unmarshal(body, &f)

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
