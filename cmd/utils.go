package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	stdinType  = "stdin"
	fileType   = "file"
	folderType = "folder"
)

const (
	csvFile = "csv"
)

// Error message constants

func getInputType(input string) (string, error) {
	if input == "" {
		return stdinType, nil
	}

	fi, err := os.Lstat(input)

	if err == nil {
		switch mode := fi.Mode(); {
		case mode.IsRegular():
			return fileType, nil
		case mode.IsDir():
			return folderType, nil
		}
	}

	return "", errors.New(fmt.Sprintf("%s: No such file or directory", input))
}

func getReaderFromInput(input string) (io.ReadCloser, error) {

	inputType, err := getInputType(input)

	if err != nil {
		return nil, err
	}

	switch inputType {
	case fileType:
		fileReader, err := os.Open(input)
		if err != nil {
			return nil, err
		}
		return fileReader, nil
	case stdinType:
		return os.Stdin, nil
	}

	return nil, err
}

func isSupportedFormat(format string, supportedFormats []string) bool {
	for _, supportedFormat := range supportedFormats {
		if format == supportedFormat {
			return true
		}
	}
	return false
}

func exitWithError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
