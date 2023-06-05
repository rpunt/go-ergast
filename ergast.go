package ergast

import (
	"encoding/json"
	"strconv"

	"github.com/rpunt/simplehttp"
)

type ErgastResults struct {
	MRData struct {
		Xmlns     string `json:"xmlns"`  // "http://ergast.com/mrd/1.5"
		Series    string `json:"series"` // "f1"
		URL       string `json:"url"`    // "http://ergast.com/api/f1/current/last/results.json"
		Limit     string `json:"limit"`  // "30"
		Offset    string `json:"offset"` // "0"
		Total     string `json:"total"`  // "20"
		RaceTable struct {
			Season string `json:"season"` // "2023"
			Round  string `json:"round"`  // "3"
			Races  []Race
		}
	}
}

type Race struct {
	Season   string `json:"season"`   // "2023"
	Round    string `json:"round"`    // "3"
	URL      string `json:"url"`      // "https://en.wikipedia.org/..."
	RaceName string `json:"raceName"` // "Australian Grand Prix"
	Circuit  Circuit
	Date     string `json:"date"` // "2023-04-02"
	Time     string `json:"time"` // "05:00:00Z"
	Results  []Result
}

type Circuit struct {
	CircuitID   string `json:"circuitId"`   // "albert_park",
	URL         string `json:"url"`         // "http://en.wikipedia.org/...",
	CircuitName string `json:"circuitName"` // "Albert Park Grand Prix Circuit"
	Location    struct {
		Lat      string `json:"lat"`      // "-37.8497"
		Long     string `json:"long"`     // "144.968"
		Locality string `json:"locality"` // "Melbourne"
		Country  string `json:"country"`  // "Australia"
	} `json:"Location"`
}

type Result struct {
	Number       string `json:"number"`
	Position     string `json:"position"`
	PositionText string `json:"positionText"`
	Points       string `json:"points"`
	Driver       Driver
	Constructor  Constructor
	Grid         string `json:"grid"`
	Laps         string `json:"laps"`
	Status       string `json:"status"`
	Time         struct {
		Millis string `json:"millis"`
		Time   string `json:"time"`
	} `json:"Time,omitempty"`
	FastestLap struct {
		Rank string `json:"rank"`
		Lap  string `json:"lap"`
		Time struct {
			Time string `json:"time"`
		} `json:"Time"`
		AverageSpeed struct {
			Units string `json:"units"`
			Speed string `json:"speed"`
		} `json:"AverageSpeed"`
	} `json:"FastestLap,omitempty"`
}

type Driver struct {
	DriverID        string `json:"driverId"`        // "max_verstappen"
	PermanentNumber string `json:"permanentNumber"` // "33"
	Code            string `json:"code"`            // "VER"
	URL             string `json:"url"`             // "http://en.wikipedia.org/..."
	GivenName       string `json:"givenName"`       // "Max"
	FamilyName      string `json:"familyName"`      // "Verstappen"
	DateOfBirth     string `json:"dateOfBirth"`     // "1997-09-30"
	Nationality     string `json:"nationality"`     // "Dutch"
}

type Constructor struct {
	ConstructorID string `json:"constructorId"` // "constructorId": "red_bull",
	URL           string `json:"url"`           // "url": "http://en.wikipedia.org/wiki/Red_Bull_Racing",
	Name          string `json:"name"`          // "name": "Red Bull",
	Nationality   string `json:"nationality"`   // "nationality": "Austrian"
}

func RaceResult() (Race, error) {
	client := simplehttp.New("https://ergast.com/api/f1")

	response, err := client.Get("/current/last/results.json")
	if err != nil {
		return Race{}, err
	}
	raceResult := ErgastResults{}
	jsonErr := json.Unmarshal([]byte(response.Body), &raceResult)
	if jsonErr != nil {
		return Race{}, jsonErr
	}

	return raceResult.MRData.RaceTable.Races[0], nil
}

// func (r Result) Status() (string, error) {
// 	for _, event := range r.SeasonContext.Timetables {
// 		if event.Description == "Race" {
// 			return event.State, nil
// 		}
// 	}

// 	return "", errors.New("unable to retrieve race timetable: no \"Race\" block")
// }

func (r Race) DriverByPosition(desiredPosition int) (Driver, error) {
	var driver Driver
	for _, result := range r.Results {
		position, err := strconv.Atoi(result.Position)
		if err != nil {
			return driver, err
		}
		if position == desiredPosition {
			driver = result.Driver
		}
	}
	return driver, nil
}

func (r Race) Winner() (Driver, error) {
	driver, err := r.DriverByPosition(1)
	if err != nil {
		return Driver{}, err
	}
	return driver, nil
}
