# owget — OpenWeatherMap CLI

A weather CLI written in Go. Supports geocoding, current weather, and 5-day forecast.

## Environment Variable

```bash
export OPENWEATHER_API_KEY="your-api-key"
```

## Build (native Go)

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o owget .
```

Build for macOS (Apple Silicon / Intel):

```bash
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o owget .
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o owget .
```

Cross compile for Linux:

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o owget .
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o owget .
```

Install to PATH:

```bash
sudo install owget /usr/local/bin/
```

Or via `go install`:

```bash
go install .
# binary goes to $GOPATH/bin/owget (make sure $GOPATH/bin is in your PATH)
```

## Build (nerdctl + Lima, multi-arch)

```bash
nerdctl.lima build --platform linux/arm64,linux/amd64 -t owget .
```

Single platform:

```bash
nerdctl.lima build -t owget .
```

## Usage

```bash
# Current weather (shortcut with coordinates)
owget 24.9575 121.5105

# Current weather
owget weather 25.0287 121.5052

# 5-day forecast
owget forecast 23.9938 120.5642

# Current weather by city name
owget city Taipei,TW

# 5-day forecast by city name
owget city Taipei,TW forecast

# Location search
owget geo Ankang,TW

# City names with spaces
owget geo "New York,US"
owget city "New York,NY,US"
owget city "New York,NY,US" forecast
```

### Flags

```bash
# Show detailed weather info (pressure, wind dir, sunrise/sunset, visibility, etc.)
owget weather Taipei,TW --detail

# Show HTTP request/response for debugging
owget geo "New York,US" --debug

# Combine both
owget city "New York,NY,US" --detail --debug
```

With nerdctl:

```bash
nerdctl.lima run --rm -e OPENWEATHER_API_KEY=$OPENWEATHER_API_KEY owget 24.9575 121.5105
nerdctl.lima run --rm -e OPENWEATHER_API_KEY=$OPENWEATHER_API_KEY owget geo "New York,US"
nerdctl.lima run --rm -e OPENWEATHER_API_KEY=$OPENWEATHER_API_KEY owget city "New York,NY,US"
```

## Coordinates Quick Reference

| Location | lat | lon |
|----------|-----|-----|
| Xindian Ankang | 24.9575 | 121.5105 |
| Taipei Zhonghua Rd | 25.0287 | 121.5052 |
| Changhua Dacun | 23.9938 | 120.5642 |
| New York, US | 40.7127 | -74.0060 |

## API Notes

Uses the OpenWeatherMap free tier:
- Geocoding: `geo/1.0/direct`
- Current Weather: `data/2.5/weather`
- 5-Day Forecast: `data/2.5/forecast`

OneCall 3.0 requires a paid subscription and is not used by this tool.
