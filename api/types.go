package api

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
