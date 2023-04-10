package ergast

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/rpunt/simplehttp"
)

type ErgastResults struct {
	MRData struct {
		Xmlns     string `json:"xmlns"`
		Series    string `json:"series"`
		URL       string `json:"url"`
		Limit     string `json:"limit"`
		Offset    string `json:"offset"`
		Total     string `json:"total"`
		RaceTable struct {
			Season string `json:"season"`
			Round  string `json:"round"`
			Races  []Race
		}
	}
}

type Race struct {
	Season   string `json:"season"`
	Round    string `json:"round"`
	URL      string `json:"url"`
	RaceName string `json:"raceName"`
	Circuit  Circuit
	Date     string `json:"date"`
	Time     string `json:"time"`
	Results  []Result
}

type Circuit struct {
	CircuitID   string `json:"circuitId"`
	URL         string `json:"url"`
	CircuitName string `json:"circuitName"`
	Location    struct {
		Lat      string `json:"lat"`
		Long     string `json:"long"`
		Locality string `json:"locality"`
		Country  string `json:"country"`
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
	DriverID        string `json:"driverId"`
	PermanentNumber string `json:"permanentNumber"`
	Code            string `json:"code"`
	URL             string `json:"url"`
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	DateOfBirth     string `json:"dateOfBirth"`
	Nationality     string `json:"nationality"`
}

type Constructor struct {
	ConstructorID string `json:"constructorId"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	Nationality   string `json:"nationality"`
}

func RaceResult() (Race, error) {
	client := simplehttp.New("https://ergast.com/api/f1")

	response, err := client.Get("/current/last/results.json")
	if err != nil {
		log.Panicf("response error: %s", err)
	}
	raceResult := ErgastResults{}
	jsonErr := json.Unmarshal([]byte(response.Body), &raceResult)
	if jsonErr != nil {
		log.Panic(jsonErr)
		// log.Fatal((jsonErr))
	}

	// return raceResult.Races[0], nil
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
		// if desiredPosition < 1 || desiredPosition > 3 {
		// 	// only positions 1-3 are tracked as RaceResults
		// 	return driver, errors.New("DriverByPosition(): valid values are 1, 2, 3")
		// }
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
