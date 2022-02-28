package junit

import (
	"encoding/xml"
)

// TestCase matches with jUnit test case
type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	ClassName string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      string   `xml:"time,attr"`
	Failure   *Failure `xml:"failure",omitempty"`
}

// TestSuites are a list of jUnit Test Suites
type TestSuites struct {
	XMLName xml.Name    `xml:"testsuites"`
	Suites  []TestSuite `xml:"testsuite"`
}

// TestSuite matches with jUnit Test Suite
type TestSuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	Tests      string     `xml:"tests,attr"`
	Failures   string     `xml:"failures,attr"`
	Time       string     `xml:"time,attr"`
	Name       string     `xml:"name,attr"`
	Properties Properties `xml:"properties"`
	TestCases  []TestCase `xml:"testcase"`
}

// Properties is a list of jUnit properties
type Properties struct {
	XMLName  xml.Name   `xml:"properties"`
	Property []Property `xml:"property"`
}

// Property is a jUnit property
type Property struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

// Failure is a jUnit failure
type Failure struct {
	XMLName        xml.Name `xml:"failure"`
	Message        string   `xml:"message,attr"`
	FailureType    string   `xml:"type,attr"`
	FailureContent string   `xml:",cdata"`
}
