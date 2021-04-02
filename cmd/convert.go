package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/allenakinkunle/swissa/converter"
)

type convertCommand struct {
	name      string
	inputFlag string
	toFlag    string
	fromFlag  string
	flagSet   *flag.FlagSet
}

func newConvertCommand() *convertCommand {
	return &convertCommand{
		name:    "convert",
		flagSet: flag.NewFlagSet("convert", flag.ExitOnError),
	}
}

func (c *convertCommand) setFlagset() {
	c.flagSet.StringVar(&c.inputFlag, "input", "", "The input you want to convert. It could either be file path or a folder path containing files you want to convert")
	c.flagSet.StringVar(&c.inputFlag, "i", "", "The input you want to convert. It could either be file path or a folder path containing files you want to convert")
	c.flagSet.StringVar(&c.fromFlag, "from", "", fmt.Sprintf("Format of the input you want to convert. Valid formats are %s", strings.Join(converter.SupportedFormats, " ")))
	c.flagSet.StringVar(&c.fromFlag, "f", "", fmt.Sprintf("Format of the input you want to convert. Valid formats are %s", strings.Join(converter.SupportedFormats, " ")))
	c.flagSet.StringVar(&c.toFlag, "to", "", fmt.Sprintf("Format you want to convert to. Valid formats are %s", strings.Join(converter.SupportedFormats, ",")))
	c.flagSet.StringVar(&c.toFlag, "t", "", fmt.Sprintf("Format you want to convert to. Valid formats are %s", strings.Join(converter.SupportedFormats, ",")))
}

func (c *convertCommand) run(args []string) {
	c.setFlagset()
	c.flagSet.Parse(args)

	if !isSupportedFormat(c.fromFlag, converter.SupportedFormats) ||
		!isSupportedFormat(c.toFlag, converter.SupportedFormats) {
		errMessage := fmt.Sprintf("Unsupported format. The valid formats are %s", strings.Join(converter.SupportedFormats, ","))
		exitWithError(errors.New(errMessage))
	}

	reader, err := getReaderFromInput(c.inputFlag)
	defer reader.Close()
	exitWithError(err)

	switch c.fromFlag {
	case csvFile:
		converter := converter.NewCSVConverter(reader)
		converter.Convert(c.toFlag, os.Stdout)
	}
}
