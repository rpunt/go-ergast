package ergast

import (
	"encoding/json"
	"encoding/xml"
	"log"

	"github.com/rpunt/simplehttp"
)

type mrdata struct {
	XMLName xml.Name `xml:"MRData"`
	Races   []Race   `xml:"RaceTable>Race"`
}

type Race struct {
	NoResults         bool
	Circuit           Circuit
	Date              ErgastDate
	Time              ErgastTime
	RaceName          string
	Season            int                `xml:"season,attr"`
	Round             int                `xml:"round,attr"`
	Results           []Result           `xml:"ResultsList>Result"`
	QualifyingResults []QualifyingResult `xml:"QualifyingList>QualifyingResult"`
}

type QualifyingResult struct {
	Driver      Driver
	Constructor Constructor
	Q1          ErgastDuration
	Q2          ErgastDuration
	Q3          ErgastDuration
	Position    int `xml:"position,attr"`
}

type Circuit struct {
	CircuitName string
}

type Result struct {
	Constructor Constructor
	Driver      Driver
	Laps        int
	Grid        int
	StatusID    int `xml:"statusId"`
	Status      string
	FastestLap  Lap
	Number      int `xml:"number,attr"`
	Position    int `xml:"position,attr"`
	Points      int `xml:"points,attr"`
}

type Lap struct {
	Time ErgastDuration
	Rank int `xml:"rank,attr"`
	Lap  int `xml:"lap,attr"`

	// There seems to be a bug with https://ergast.com/api/f1/2017/20/results which makes the averagespeed field match the time field.
	// Leaving this field unimplemented for now. A bug report has been submitted to the ergast API
	//AverageSpeed float64
	//AverageSpeedUnits string `xml:"units,attr"`
}

type Driver struct {
	DriverID        string `xml:"driverId,attr"`
	Code            string `xml:"code,attr"`
	PermanentNumber int
	GivenName       string
	FamilyName      string
	Nationality     string
	DateOfBirth     ErgastDate
}

type Constructor struct {
	ConstructorID string `xml:"constructorId,attr"`
	Name          string
	Nationality   string
}

func RaceResult() (Race, error) {
	client := simplehttp.New()
	client.BaseURL = "https://ergast.com/api/f1"

	response, err := client.Get("/current/last/results")
	if err != nil {
		log.Panicf("response error: %s", err)
	}
	raceResult := mrdata{}
	jsonErr := json.Unmarshal([]byte(response.Body), &raceResult)
	if jsonErr != nil {
		log.Panic(jsonErr)
	}

	return raceResult.Races[0], nil
}
