package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	junit "github.com/RH-QE-Distributed-Tracing/junitcli/pkg/model"
)

const (
	flagVerbose    = "verbose"
	flagSuiteName  = "suite-name"
	flagOutput     = "output"
	flagShowReport = "report"
)

func readXML(xmlPathName string) (junit.TestSuites, error) {
	bytes, err := ioutil.ReadFile(xmlPathName)
	if err != nil {
		return junit.TestSuites{}, err
	}

	var suites junit.TestSuites
	err = xml.Unmarshal(bytes, &suites)

	if err != nil {
		return junit.TestSuites{}, err
	}

	return suites, nil
}

func initCmd() error {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetDefault(flagVerbose, false)
	flag.Bool(flagVerbose, false, "Enable verbose output")

	viper.SetDefault(flagShowReport, false)
	flag.Bool(flagShowReport, false, "Show table report")

	viper.SetDefault(flagSuiteName, "")
	flag.String(flagSuiteName, "", "Suite name to set (just applies if there is one suite)")

	viper.SetDefault(flagOutput, "")
	flag.String(flagOutput, "", "Output file")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	flag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err == nil {
		if viper.GetBool(flagVerbose) {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	return err
}

func main() {
	err := initCmd()
	if err != nil {
		logrus.Fatalln(err)
	}

	xmlPathName := ""

	switch flag.NArg() {
	default:
		flag.Usage()
		os.Exit(1)
	case 1:
		xmlPathName = flag.Args()[0]
	}

	fileInfo, err := os.Stat(xmlPathName)
	if err != nil {
		logrus.Fatalln("Error while opening the jUnit reports: ", err)
	}

	var suites junit.TestSuites

	if fileInfo.IsDir() {
		err = filepath.Walk(xmlPathName,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !strings.HasSuffix(path, ".xml") {
					return nil
				}

				newSuites, err := readXML(path)
				if err != nil {
					return err
				}

				suites.Aggregate(&newSuites)

				return nil
			})
		if err != nil {
			logrus.Fatalln(err)
		}
	} else {
		suites, err = readXML(xmlPathName)
		if err != nil {
			logrus.Fatalln("Error while reading the file ", xmlPathName, ":", err)
		}

		// Change test suite name
		if viper.GetString(flagSuiteName) != "" {
			err = suites.SetTestSuiteName(viper.GetString(flagSuiteName))
			if err != nil {
				logrus.Fatalln(err)
			}
		}
	}
	suites.Sanitize()

	// Show report
	if viper.GetBool(flagShowReport) {
		err = suites.DrawReport()
		if err != nil {
			logrus.Fatalln(err)
		}
	}

	outputFile := viper.GetString(flagOutput)
	if outputFile != "" {
		outputFileContent, err := xml.MarshalIndent(suites, "", "    ")
		if err != nil {
			logrus.Fatalln("There was a problem while creating the output file", err)
		}

		err = os.WriteFile(outputFile, outputFileContent, 0644)
	}

}
