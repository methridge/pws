package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fatih/color"
)

type weatherCurrent struct {
	Observations []struct {
		StationID         string      `json:"stationID"`
		ObsTimeUtc        time.Time   `json:"obsTimeUtc"`
		ObsTimeLocal      string      `json:"obsTimeLocal"`
		Neighborhood      string      `json:"neighborhood"`
		SoftwareType      string      `json:"softwareType"`
		Country           string      `json:"country"`
		SolarRadiation    float64     `json:"solarRadiation"`
		Lon               float64     `json:"lon"`
		RealtimeFrequency interface{} `json:"realtimeFrequency"`
		Epoch             int         `json:"epoch"`
		Lat               float64     `json:"lat"`
		Uv                float64     `json:"uv"`
		Winddir           int         `json:"winddir"`
		Humidity          int         `json:"humidity"`
		QcStatus          int         `json:"qcStatus"`
		Imperial          struct {
			Temp        int     `json:"temp"`
			HeatIndex   int     `json:"heatIndex"`
			Dewpt       int     `json:"dewpt"`
			WindChill   int     `json:"windChill"`
			WindSpeed   int     `json:"windSpeed"`
			WindGust    int     `json:"windGust"`
			Pressure    float64 `json:"pressure"`
			PrecipRate  float64 `json:"precipRate"`
			PrecipTotal float64 `json:"precipTotal"`
			Elev        int     `json:"elev"`
		} `json:"imperial"`
	} `json:"observations"`
}

var (
	api   string
	sid   string
	units string
	key   string
)

func main() {
	url := api + "?stationId=" + sid + "&format=json&units=" + units + "&apiKey=" + key
	// fmt.Println("URL: ", url)
	// fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject weatherCurrent
	json.Unmarshal(bodyBytes, &responseObject)
	compassDirs := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW", "N"}
	c := color.New(color.FgCyan)
	g := color.New(color.FgGreen)
	y := color.New(color.FgYellow)

	// Header
	c.Printf("Current Conditions for ")
	y.Printf("%s", responseObject.Observations[0].StationID)
	c.Printf(" at ")
	y.Printf("%s", responseObject.Observations[0].ObsTimeLocal)
	c.Printf(" are:\n")

	// Current
	c.Printf("Current:    ")
	if responseObject.Observations[0].Imperial.Temp > 80 {
		color.Set(color.FgRed)
	} else if 60 < responseObject.Observations[0].Imperial.Temp && responseObject.Observations[0].Imperial.Temp < 80 {
		color.Set(color.FgGreen)
	} else {
		color.Set(color.FgHiBlue)
	}
	fmt.Printf("%d\u00B0F (%d\u00B0C)\n",
		responseObject.Observations[0].Imperial.Temp,
		(((responseObject.Observations[0].Imperial.Temp - 32) * 5) / 9),
	)

	// Feels Like
	c.Printf("Feels Like: ")
	if responseObject.Observations[0].Imperial.Temp > 80 {
		color.Set(color.FgRed)
	} else if 60 < responseObject.Observations[0].Imperial.Temp && responseObject.Observations[0].Imperial.Temp < 80 {
		color.Set(color.FgGreen)
	} else {
		color.Set(color.FgHiBlue)
	}
	fmt.Printf("%d\u00B0F (%d\u00B0C)\n",
		responseObject.Observations[0].Imperial.WindChill,
		(((responseObject.Observations[0].Imperial.WindChill - 32) * 5) / 9),
	)

	// Dew Point
	c.Printf("Dew Point:  ")
	g.Printf("%d\u00B0F (%d\u00B0C)\n",
		responseObject.Observations[0].Imperial.Dewpt,
		(((responseObject.Observations[0].Imperial.Dewpt - 32) * 5) / 9),
	)

	// Humidity
	c.Printf("Humidity:   ")
	g.Printf("%d%%\n",
		responseObject.Observations[0].Humidity,
	)

	// Wind Direction
	compassIndex := responseObject.Observations[0].Winddir / 22
	c.Printf("Wind:       ")
	g.Printf("%s(%d\u00B0) @ %d-%d mph\n",
		compassDirs[compassIndex],
		responseObject.Observations[0].Winddir,
		responseObject.Observations[0].Imperial.WindSpeed,
		responseObject.Observations[0].Imperial.WindGust,
	)
}
