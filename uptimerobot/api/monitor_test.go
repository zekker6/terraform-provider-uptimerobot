package uptimerobotapi

import (
	"testing"
)

func TestMonitor_ParseCustomHttpStatuses(t *testing.T) {
	badInputs := []string{
		"404:6",   // invalid type
		"asd:0",   // invalid status
		"asd:asd", // invalid status and type
		":1",      // invalid status
		"200:",    // invalid type
	}

	for _, input := range badInputs {
		err, _, _ := ParseCustomHttpStatuses(input)
		if err == nil {
			t.Error("Parsing should return error", input)
		}
	}

	err, successCodes, downCodes := ParseCustomHttpStatuses("404:0_200:1")
	if err != nil {
		t.Error("Parsing failed", err)
	}

	if len(downCodes) != 1 {
		t.Error("Expected to have single item parsed")
	}

	if len(successCodes) != 1 {
		t.Error("Expected to have single item parsed")
	}

	if len(downCodes) == 1 && downCodes[0] != 404 {
		t.Error("Parsed item doesnt match expected")
	}

	err, downCodes, downCodes = ParseCustomHttpStatuses("")
	if err != nil {
		t.Error("Should not return error", err)
	}
	if len(downCodes) != 0 {
		t.Errorf("Expected to have no items parsed")
	}

	if len(downCodes) != 0 {
		t.Errorf("Expected to have no items parsed")
	}
}

func TestMonitor_BuildCustomHttpStatusString(t *testing.T) {
	successCodes := []int{
		200,
		301,
	}

	errorCodes := []int{
		404,
		500,
	}

	res := BuildCustomHttpStatusString(successCodes, errorCodes)
	if res != "200:1_301:1_404:0_500:0" {
		t.Error("Build not matching expected result", successCodes, errorCodes, res)
	}

	res = BuildCustomHttpStatusString([]int{}, []int{})
	if res != "" {
		t.Error("Build not matching expected result", successCodes, errorCodes, res)
	}
}

func TestMonitor_CustomHttpStatusStringParsersCompatibility(t *testing.T) {
	successCodes := []int{
		200,
		301,
	}

	errorCodes := []int{
		404,
		500,
	}
	res := BuildCustomHttpStatusString(successCodes, errorCodes)

	_, parsedSuccess, parsedDown := ParseCustomHttpStatuses(res)

	arraysEqual := func(one, other []int) bool {
		if len(one) != len(other) {
			return false
		}

		for keyOne, valOne := range one {
			valOther := other[keyOne]
			if valOne != valOther {
				return false
			}
		}

		return true
	}

	if !arraysEqual(parsedDown, errorCodes) {
		t.Error("Parse results validation failed", parsedDown, errorCodes)
	}

	if !arraysEqual(parsedSuccess, successCodes) {
		t.Error("Parse results validation failed", parsedSuccess, successCodes)
	}

}
