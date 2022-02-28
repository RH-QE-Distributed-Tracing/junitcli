package main

import (
	"flag"
	"os"
	"strings"

	junit "github.com/joshdk/go-junit"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagVerbose = "verbose"
)

func readXML(xmlFileName string) ([]junit.Suite, error) {
	logrus.Debugln("Reading file", xmlFileName)
	return junit.IngestFile(xmlFileName)
}

func drawTable(suites []junit.Suite) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Result"})

	for _, testSuite := range suites {
		for _, testCase := range testSuite.Tests {
			table.Append([]string{testCase.Name, string(testCase.Status)})
		}
	}
	table.Render()
}

func initCmd() error {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetDefault(flagVerbose, false)
	flag.Bool(flagVerbose, false, "Enable verbose output")

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

	xmlFileName := ""

	switch flag.NArg() {
	default:
		flag.Usage()
		os.Exit(1)
	case 1:
		xmlFileName = flag.Args()[0]
	}

	suites, err := readXML(xmlFileName)

	if err != nil {
		logrus.Fatalln(err)
	}

	drawTable(suites)
}
