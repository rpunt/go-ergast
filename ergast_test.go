package ergast

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestErgast(t *testing.T) {
	race, err := RaceResult()

	if err != nil {
		t.Errorf(`RaceResult("") = %q, %v, want "", error`, race, err)
	}
}

func TestDriverByPosition(t *testing.T) {
	response := ReadTestData("test_data/verstappen_wins.json")
	race := response.MRData.RaceTable.Races[0]

	position := 3
	want := "ALO"
	if got, _ := race.DriverByPosition(position); got.Code != want {
		t.Errorf(`DriverByPosition(%v) = %q, want "%v"`, position, got.Code, want)
	}

	// only driver positions 1-3 are returned in this API endpoint
	// test bounds checking
	// for position := 0; position < 5; position += 4 {
	// 	if _, err := race.DriverByPosition(position); err == nil {
	// 		t.Errorf(`DriverByPosition(%v) == nil, expected bounds-checking error`, position)
	// 	}
	// }
}

func TestResultParsing(t *testing.T) {
	response := ReadTestData("test_data/verstappen_wins.json")
	race := response.MRData.RaceTable.Races[0]
	want := "VER"
	got, _ := race.Winner()
	if got.Code != want {
		t.Errorf(`Winner() = %q, want "%v"`, got, want)
	}

	// want = "completed"
	// if got, _ := race.Status(); got != want {
	// 	t.Errorf(`Status() = %q, want "%v"`, got, want)
	// }
}

func ReadTestData(filename string) ErgastResults {
	content, err := os.ReadFile(filename)
	CheckErr(err)

	race := ErgastResults{}
	jsonErr := json.Unmarshal(content, &race)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return race
}

// Check for and log errors
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
