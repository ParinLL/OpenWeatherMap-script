package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	GeoURL      = "https://api.openweathermap.org/geo/1.0/direct"
	WeatherURL  = "https://api.openweathermap.org/data/2.5/weather"
	ForecastURL = "https://api.openweathermap.org/data/2.5/forecast"
)

var Debug bool

// EncodeQuery URL-encodes each comma-separated part of the query
// while preserving commas as-is (OpenWeatherMap uses commas as delimiters).
func EncodeQuery(query string) string {
	parts := strings.Split(query, ",")
	for i, p := range parts {
		parts[i] = url.QueryEscape(strings.TrimSpace(p))
	}
	return strings.Join(parts, ",")
}

func HTTPGet(reqURL string) []byte {
	if Debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] GET %s\n", RedactURLCredentials(reqURL))
	}
	resp, err := http.Get(reqURL)
	if err != nil {
		Fatal("http error: " + err.Error())
	}
	defer resp.Body.Close()
	if Debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] HTTP %d\n", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Fatal("read error: " + err.Error())
	}
	if Debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Body: %s\n", string(body))
	}
	if resp.StatusCode != 200 {
		Fatal(fmt.Sprintf("API %d: %s", resp.StatusCode, string(body)))
	}
	return body
}

// RedactURLCredentials masks sensitive query params in debug logs.
func RedactURLCredentials(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	q := u.Query()
	for _, key := range []string{"appid", "api_key", "apikey", "token", "access_token"} {
		if _, ok := q[key]; ok {
			q.Set(key, "[REDACTED]")
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func Unmarshal(data []byte, v any) {
	if err := json.Unmarshal(data, v); err != nil {
		Fatal("json error: " + err.Error())
	}
}

func Fatal(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}
