package main

import (
	"Go-BticinoThermostatInflux/config"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

// Devices Array of thermostats status
type Devices struct {
	Devices []DeviceStatus `json:"chronothermostats"`
}

// Measurement The struct of the measurements in the DeviceStatus struct
type Measurement struct {
	TimeStamp string `json:"timeStamp"`
	Value     string `json:"value"`
	Unit      string `json:"unit"`
}

// Measurements Just an array (the measurements are an array for unknown reasons)
type Measurements struct {
	Measures []Measurement `json:"measures"`
}

// DeviceStatus The status of a Thermostat
type DeviceStatus struct {
	Function          string          `json:"function"`
	Mode              string          `json:"mode"`
	Setpoint          json.RawMessage `json:"setpoint"`
	Programs          json.RawMessage `json:"programs"`
	TemperatureFormat string          `json:"temperatureformat"`
	LoadState         string          `json:"loadstate"`
	Time              time.Time       `json:"time"`
	Thermometer       Measurements    `json:"thermometer"`
	Hygrometer        Measurements    `json:"hygrometer"`
	Sender            json.RawMessage `json:"sender"`
}

// GetToken The response of getRefreshToken and getAccessToken
type GetToken struct {
	AccessToken           string `json:"access_token"`
	IDToken               string `json:"id_token"`
	TokenType             string `json:"token_type"`
	NotBefore             string `json:"not_before"`
	ExpiresIn             uint16 `json:"expires_in"`
	ExpiresOn             string `json:"expires_on"`
	Resource              string `json:"resource"`
	IDTokenExpiresIn      uint16 `json:"id_token_expires_in"`
	ProfileInfo           string `json:"profile_info"`
	Scope                 string `json:"scope"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn uint16 `json:"refresh_token_expires_in"`
}

func sigtermHandler(influx influxdb2.Client) {
	// Prepare to catch the SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Got SIGTERM!")
		// Close InfluxDB Connection
		influx.Close()
		os.Exit(0)
	}()
}

func getAuthToken() (authToken string) {
	fmt.Println("Open the following link in a browser and login. Once logged in, it'll redirect you to google.com?code=<something>. Type in that something")
	fmt.Println(config.AuthEndpoint + "authorize?client_id=" + config.ClientID + "&response_type=code&redirect_uri=google.com")

	// Read the console input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter code: ")
	code, _ := reader.ReadString('\n')

	return strings.TrimSpace(code)
}

func getRefreshCode(accessToken string) (refreshToken string) {
	// Create POST request payload
	data := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {accessToken},
	}
	resp, err := http.PostForm(config.AuthEndpoint+"token", data)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Elaborate resp
	body, _ := ioutil.ReadAll(resp.Body)

	// Initialize the variable that'll contain all the data
	var token GetToken

	// Parse the JSON body
	json.Unmarshal(body, &token)
	refreshToken = token.RefreshToken

	fmt.Println(string(body))

	// Write it in the file
	err = ioutil.WriteFile(config.RefreshFileName, []byte(refreshToken), 0700)

	// Handle errors while writing the file
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}

	return refreshToken
}

func refreshTokenFlow(refreshToken string) (updatedRefreshToken string, accessToken string) {
	// Create POST request payload
	data := url.Values{
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	resp, err := http.PostForm(config.AuthEndpoint+"token", data)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Elaborate resp
	body, _ := ioutil.ReadAll(resp.Body)

	// Initialize the variable that'll contain all the data
	var token GetToken

	// Parse the JSON body
	json.Unmarshal(body, &token)
	updatedRefreshToken = token.RefreshToken
	accessToken = token.AccessToken

	// Write it in the file
	err = ioutil.WriteFile(config.RefreshFileName, []byte(updatedRefreshToken), 0700)

	// Handle errors while writing the file
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}

	return updatedRefreshToken, accessToken
}

func auth() (accessToken string) {
	if _, err := os.Stat(config.RefreshFileName); err == nil {
		// If there's a refresh.txt file, try to use that refresh token
		fileData, err := ioutil.ReadFile(config.RefreshFileName)

		// Handle eventual error
		if err != nil {
			panic("Unable to read file")
		}
		refreshToken := strings.TrimSpace(string(fileData))

		refreshToken, accessToken = refreshTokenFlow(refreshToken)
		if refreshToken == "" || accessToken == "" {
			accessToken = getAuthToken()
			getRefreshCode(accessToken)
		}
	} else {
		accessToken = getAuthToken()
		getRefreshCode(accessToken)
	}
	return accessToken
}

func getThermostatStatus(accessToken string, plantID string, moduleID string) (temperature string, humidity string, status bool) {
	// Generate URL and make the GET request
	url := config.APIEndpoint + "chronothermostat/thermoregulation/addressLocation/plants/" + plantID + "/modules/parameter/id/value/" + moduleID

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Ocp-Apim-Subscription-Key", config.SubscriptionKey)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)

	// Handle the eventual error
	if err != nil {
		panic(err)
	}

	// Close the response body (Why? Because the docs say so)
	defer resp.Body.Close()

	// If no error is found, get the request body and parse it
	byteValue, _ := ioutil.ReadAll(resp.Body)

	// Initialize the variable that'll contain all the data
	var thermostat Devices

	// Parse the JSON body
	err = json.Unmarshal(byteValue, &thermostat)
	if err != nil {
		panic(err)
	}

	// Extract the needed data from the struct
	temperature = thermostat.Devices[0].Thermometer.Measures[0].Value
	humidity = thermostat.Devices[0].Hygrometer.Measures[0].Value

	if thermostat.Devices[0].LoadState == "ACTIVE" { // When the LoadState is "ACTIVE", then the thermostat is heating
		status = true
	} else {
		status = false
	}

	return temperature, humidity, status
}

func main() {
	// Create a InfluxDB Client to push the data in the DB
	client := influxdb2.NewClient(config.InfluxHost, config.InfluxToken)
	writeAPI := client.WriteAPI(config.InfluxOrg, config.InfluxBucket)

	// Start the SIGTERM handler
	sigtermHandler(client)

	for true {
		// Authenticate and get data
		accessToken := auth()
		temperature, humidity, status := getThermostatStatus(accessToken, config.PlantID, config.ModuleID)

		// Get the current time (to add to the data)
		relTime := time.Now()

		// Create the point with all the data
		p := influxdb2.NewPointWithMeasurement("environment").
			AddTag("sensorType", "BticinoThermostat").
			AddTag("sensorID", config.ModuleID).
			AddField("temperature", temperature).
			AddField("humidity", humidity).
			AddField("status", status).
			SetTime(relTime)
		writeAPI.WritePoint(p)
		writeAPI.Flush()

		time.Sleep(15 * time.Second)
	}
}
