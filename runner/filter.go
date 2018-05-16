package runner

import (
	"encoding/json"
	"fmt"
)

type TestCaseFilter interface {
	Filter([]TestCase) ([]TestCase, error)
}

type IntegrationConfigTestCaseFilter map[string]interface{}

func NewIntegrationConfigTestCaseFilter(rawConfig []byte) (IntegrationConfigTestCaseFilter, error) {
	filter := IntegrationConfigTestCaseFilter{}
	err := json.Unmarshal(rawConfig, &filter)
	return filter, err
}

func (f IntegrationConfigTestCaseFilter) Filter(testCases []TestCase) ([]TestCase, error) {
	var filteredTestCases []TestCase
	for _, testCase := range testCases {
		flagValue, err := f.getFlagValue(testCase.Name())
		if err != nil {
			return nil, err
		}

		if flagValue == true {
			filteredTestCases = append(filteredTestCases, testCase)
		}
	}

	if (len(filteredTestCases)) > 0 {
		return filteredTestCases, nil
	}

	return nil, fmt.Errorf("unable to find any test case included by the config")
}

func (f IntegrationConfigTestCaseFilter) getFlagValue(testCaseName string) (bool, error) {
	flagName := fmt.Sprintf("include_%s", testCaseName)

	flagValue, isDefined := f[flagName]
	if !isDefined {
		flagValue = false
	}

	boolFlagValue, isBool := flagValue.(bool)
	if !isBool {
		return false, fmt.Errorf("'%s' should be a boolean", flagName)
	}

	return boolFlagValue, nil
}
