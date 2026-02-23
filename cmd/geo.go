package cmd

import (
	"fmt"
	"strings"

	"owget/api"
)

func Geo(apiKey, query string, detail bool) {
	reqURL := fmt.Sprintf("%s?q=%s&limit=5&appid=%s", api.GeoURL, api.EncodeQuery(query), apiKey)
	body := api.HTTPGet(reqURL)

	var results []api.GeoResult
	api.Unmarshal(body, &results)

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

func City(apiKey, query string, extra []string, detail bool) {
	reqURL := fmt.Sprintf("%s?q=%s&limit=1&appid=%s", api.GeoURL, api.EncodeQuery(query), apiKey)
	body := api.HTTPGet(reqURL)

	var results []api.GeoResult
	api.Unmarshal(body, &results)

	if len(results) == 0 {
		api.Fatal("city not found: " + query)
	}

	r := results[0]
	fmt.Printf("📍 %s, %s  (lat=%.4f, lon=%.4f)\n\n", r.Name, r.Country, r.Lat, r.Lon)

	if len(extra) > 0 && extra[0] == "forecast" {
		Forecast(apiKey, r.Lat, r.Lon, detail)
	} else {
		Weather(apiKey, r.Lat, r.Lon, detail)
	}
}
