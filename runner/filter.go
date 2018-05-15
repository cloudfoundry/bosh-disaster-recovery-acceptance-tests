package runner

import (
	"fmt"
)

type TestCaseFilter interface {
	Filter([]TestCase) []TestCase
}

type IntegrationConfigTestCaseFilter map[string]interface{}

func (f IntegrationConfigTestCaseFilter) Filter(testCases []TestCase) []TestCase {
	var filteredTestCases []TestCase
	for _, testCase := range testCases {
		if f["include_"+testCase.Name()] == true {
			filteredTestCases = append(filteredTestCases, testCase)
		}
	}

	if (len(filteredTestCases)) > 0 {
		return filteredTestCases
	}

	panic(fmt.Sprintf("Unable to find any test case included by the config"))
}
