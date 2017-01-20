package main

import (
	"bytes"
	"fmt"
	"gocsv"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Sensor struct {
	Lat         string `csv:"lat"`
	Lon         string `csv:"lon"`
	Stationcode string `csv:"stationCode"`
	State       string `csv:"state"`
	Title       string `csv:"title"`
	RawDate     string `csv:"date"`
	Date        string
	RawVal      string `csv:"val"`
	Val         string
	Type2       string `csv:"type"`
	Valtype     string `csv:"valtype"`
	Unit        string
}

func main() {

	// Currentdate
	t := time.Now()
	currentDate, _ := strconv.Atoi(t.Format("20060102"))

	// Inserting Data for last 14 Days (last day is day prior)
	for i := 0; i < 14; i++ {
		currentDate -= 1

		// Creating Sensor Objects
		Data := getData(currentDate)

		// Format Values inside Sensor struct
		formatValues(Data)

		// Filling Templates and posting Data to SOS
		// Insert Sensors only once
		if i == 0 {
			insertDataIntoSOS(Data, "InsertSensor.xml")
		}
		insertDataIntoSOS(Data, "InsertObservation.xml")
	}
}

func getData(currentDate int) []*Sensor {

	// Filling currentDate into URL
	url := fmt.Sprintf("https://www.umweltbundesamt.de/luftdaten/stations/locations?pollutant=PM1&data_type=1TMW&date=%v", currentDate)

	// Build HTTP GET Request
	req, err := http.NewRequest("GET", url, nil)
	// Fehlerbehandlung
	if err != nil {
		log.Fatal("Error on building Request to Data Server", err)
		return nil
	}
	req.Close = true

	client := &http.Client{}
	resp, err := client.Do(req)
	// Fehlerbehandlung
	if err != nil {
		log.Fatal("Error on Response from Data Server", err)
		return nil
	}
	// No idea what this is for
	defer resp.Body.Close()

	// Reading Response
	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading Data", err)
		return nil
	}

	// Parsing Response to proper CSV Format
	responseBodyCSV := string(responseBodyBytes)
	responseBodyCSV = strings.Replace(responseBodyCSV, ",", " ", -1)
	responseBodyCSV = strings.Replace(responseBodyCSV, "	", "\",\"", -1)
	responseBodyCSV = strings.Replace(responseBodyCSV, "\n", "\"\n\"", -1)
	responseBodyCSV = strings.Replace(responseBodyCSV, "&uuml;", "ü", -1)
	responseBodyCSV = strings.Replace(responseBodyCSV, "&auml;", "ä", -1)
	responseBodyCSV = strings.Replace(responseBodyCSV, "&ouml;", "ö", -1)
	responseBodyCSV = strings.Replace(responseBodyCSV, "&szlig;", "ß", -1)
	responseBodyCSV = "\"" + responseBodyCSV + "\""

	// Slice for Storing Sensors
	sensors := []*Sensor{}

	// CSV to struct
	if err := gocsv.UnmarshalString(responseBodyCSV, &sensors); err != nil {
		log.Fatal("Error Converting CSV to struct")
	}
	resp.Body.Close()

	return sensors
}
func insertDataIntoSOS(sensors []*Sensor, filledtemplate string) bool {

	// Loading Template
	var template = template.Must(template.ParseFiles("./template/" + filledtemplate))

	url := "http://127.0.0.1:9001/52n-sos-webapp/service"

	for _, sensor := range sensors {
		// Buffer Filled Template
		var buffer bytes.Buffer
		// Fill Template with Data
		if err := template.Execute(&buffer, sensor); err != nil {
			log.Fatal("Error on building Template for Request to Data Server", err)
			return false
		}
		postingDataToSOS(url, buffer)
	}
	return true
}

func postingDataToSOS(url string, body bytes.Buffer) bool {

	// Creating POST Request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body.Bytes()))
	if err != nil {
		log.Fatal("Error on creating POST Request")
		return false
	}

	// Setting appropiate Content Header
	req.Header.Set("Content-Type", "application/xml")

	// Posting POST Request to Server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer resp.Body.Close()

	// Uncomment to Export POST Responses to Commandline
	/*
		responseBodyBytes2, err := ioutil.ReadAll(resp.Body)
		fmt.Println(bytes.NewBuffer(responseBodyBytes2).String()) */

	resp.Body.Close()
	return true
}

func formatValues(sensors []*Sensor) {
	for _, sensor := range sensors {
		date := strings.Split(sensor.RawDate, ".")
		sensor.Date = date[2] + "-" + date[1] + "-" + date[0]

		value := strings.Split(sensor.RawVal, " ")
		sensor.Val = value[0]
		sensor.Unit = value[1]
	}
}
