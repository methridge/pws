package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const reset = "\033[0m"
const red = "\033[0;31m"
const green = "\033[0;32m"
const yellow = "\033[0;33m"
const blue = "\033[0;34m"
const purple = "\033[0;35m"
const cyan = "\033[0;36m"
const gray = "\033[0;37m"
const white = "\033[97m"

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

func main() {
	viper.SetConfigName("pws")
	viper.SetConfigType("hcl")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	url := viper.GetString("api") + "?stationId=" + viper.GetString("sid") + "&format=json&units=" + viper.GetString("units") + "&apiKey=" + viper.GetString("key")
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
	fmt.Printf("Current Conditions for %s%s%s at %s%s%s are:\n",
		yellow,
		responseObject.Observations[0].StationID,
		reset,
		yellow,
		responseObject.Observations[0].ObsTimeLocal,
		reset,
	)
	fmt.Printf("%sCurrent:%s    %s%d\u00B0F (%d\u00B0C)%s\n",
		cyan,
		reset,
		green,
		responseObject.Observations[0].Imperial.Temp,
		(((responseObject.Observations[0].Imperial.Temp - 32) * 5) / 9),
		reset,
	)
	if responseObject.Observations[0].Imperial.Temp > 70 {
		fmt.Printf("%sFeels Like:%s %s%d\u00B0F (%d\u00B0C)%s\n",
			cyan,
			reset,
			green,
			responseObject.Observations[0].Imperial.HeatIndex,
			(((responseObject.Observations[0].Imperial.HeatIndex - 32) * 5) / 9),
			reset,
		)
	} else {
		fmt.Printf("%sFeels Like:%s %s%d\u00B0F (%d\u00B0C)%s\n",
			cyan,
			reset,
			green,
			responseObject.Observations[0].Imperial.WindChill,
			(((responseObject.Observations[0].Imperial.WindChill - 32) * 5) / 9),
			reset,
		)
	}
	fmt.Printf("%sDew Point:%s  %s%d\u00B0F (%d\u00B0C)%s\n",
		cyan,
		reset,
		green,
		responseObject.Observations[0].Imperial.Dewpt,
		(((responseObject.Observations[0].Imperial.Dewpt - 32) * 5) / 9),
		reset,
	)
	fmt.Printf("%sHumidity:%s   %s%d%%%s\n",
		cyan,
		reset,
		green,
		responseObject.Observations[0].Humidity,
		reset,
	)
	compassIndex := responseObject.Observations[0].Winddir / 22
	fmt.Printf("%sWind:%s       %s%s(%d\u00B0) @ %d-%d mph%s\n",
		cyan,
		reset,
		green,
		compassDirs[compassIndex],
		responseObject.Observations[0].Winddir,
		responseObject.Observations[0].Imperial.WindSpeed,
		responseObject.Observations[0].Imperial.WindGust,
		reset,
	)
}
