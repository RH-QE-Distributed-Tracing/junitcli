package junit

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

// Failure from jUnit report
type Failure struct {
	XMLName        xml.Name `xml:"failure"`
	Message        string   `xml:"message,attr"`
	FailureType    string   `xml:"type,attr"`
	FailureContent string   `xml:",cdata"`
}

// TestCase from jUnit report
type TestCase struct {
	XMLName    xml.Name `xml:"testcase"`
	ClassName  string   `xml:"classname,attr"`
	Name       string   `xml:"name,attr"`
	Time       float64  `xml:"time,attr"`
	Failure    *Failure `xml:"failure,omitempty"`
	Assertions int      `xml:"assertions,attr"`
}

// TestSuites from jUnit report
type TestSuites struct {
	XMLName  xml.Name    `xml:"testsuites"`
	Name     string      `xml:"name,attr"`
	Failures int         `xml:"failures,attr"`
	Time     float64     `xml:"time,attr"`
	Tests    int         `xml:"tests,attr"`
	Suites   []TestSuite `xml:"testsuite"`
}

// TestSuite from jUnit report
type TestSuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Tests     int        `xml:"tests,attr"`
	Failures  int        `xml:"failures,attr"`
	Time      float64    `xml:"time,attr"`
	Name      string     `xml:"name,attr"`
	TestCases []TestCase `xml:"testcase"`
}

// normalizeNames uses names that can be ingested by any jUnit report tool.
func (suites *TestSuites) normalizeNames() {
	logrus.Debugln("normalizing test case names")

	if len(suites.Suites) == 0 {
		logrus.Warningln("No error suites found! No report printed")
		return
	}

	regex, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		logrus.Fatal(err)
	}

	for i := 0; i < len(suites.Suites); i++ {
		for j := 0; j < len(suites.Suites[i].TestCases); j++ {
			name := suites.Suites[i].TestCases[j].Name
			name = strings.TrimSpace(name)
			name = strings.ReplaceAll(name, " ", "_")
			name = strings.ReplaceAll(name, "-", "_")
			// Remove non alpha numeric characters
			name = regex.ReplaceAllString(name, "")

			logrus.Debugln(suites.Suites[i].TestCases[j].Name, "->", name)

			suites.Suites[i].TestCases[j].Name = name
			// Set a value to ensure the parsers can read it
			suites.Suites[i].TestCases[j].ClassName ="test"
		}
	}

}

// pruneArtifactsTestCase removes the artifacts TestCase
func (suites *TestSuites) pruneArtifactsTestCase() {
	logrus.Debugln("Removing 'artifacts' TestCase")

	var artifactsIndex int = -1
	var totalTests int = 0
	for i := 0; i < len(suites.Suites); i++ {
		for j := 0; j < len(suites.Suites[i].TestCases); j++ {
			if suites.Suites[i].TestCases[j].Name == "artifacts" {
				artifactsIndex = j
				break
			}
		}

		if artifactsIndex != -1 {
			testCaseList := append(
				suites.Suites[i].TestCases[:artifactsIndex],
				suites.Suites[i].TestCases[artifactsIndex+1:]...,
			)

			suites.Suites[i].TestCases = testCaseList
			suites.Suites[i].Tests = len(testCaseList)

			artifactsIndex = -1
			logrus.Debugln("Test case 'artifacts' removed!")
		}
		totalTests += suites.Suites[i].Tests
	}
	suites.Tests = totalTests
}

// Sanitize will clean the name of some of the tests and remove the artifacts
// test (which is created by KUTTL and does nothing)
func (suites *TestSuites) Sanitize() {
	suites.pruneArtifactsTestCase()
	suites.normalizeNames()
}

// DrawReport prints a report in the console
func (suites *TestSuites) DrawReport() error {

	if len(suites.Suites) == 0 {
		return fmt.Errorf("no error suites found! No report printed")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Result"})

	for _, testSuite := range suites.Suites {
		for _, testCase := range testSuite.TestCases {
			var status string
			if testCase.IsPassed() {
				status = "passed"
			} else {
				status = "failed"
			}
			table.Append([]string{testCase.Name, status})
		}
	}
	table.Render()
	return nil
}

// SetTestSuiteName sets the name for a single test suite. This is needed because KUTTL
// doesn't set the name of the suites
func (suites *TestSuites) SetTestSuiteName(newSuiteName string) error {
	logrus.Debugln("Setting a new name for the test suites")
	if len(suites.Suites) == 0 {
		return fmt.Errorf("no test suites found")
	} else if len(suites.Suites) > 1 {
		return fmt.Errorf("more than 1 suite found. SetName method can not be used")
	}

	// Change the name to the suite
	suites.Suites[0].Name = newSuiteName

	for i := 0; i < len(suites.Suites[0].TestCases); i++ {
		testName := suites.Suites[0].TestCases[i].Name
		suites.Suites[0].TestCases[i].Name = fmt.Sprintf("%s/%s", newSuiteName, testName)
	}

	return nil
}

// Aggregate includes another suites under this one
func (suites *TestSuites) Aggregate(anotherSuite *TestSuites) {
	suites.Suites = append(suites.Suites, anotherSuite.Suites...)
}

// IsPassed returns true if there were not error with the tests
func (test *TestCase) IsPassed() bool {
	return test.Failure == nil
}
