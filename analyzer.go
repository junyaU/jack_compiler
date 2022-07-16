package jack_compiler

import (
	"io/ioutil"
	"os"
)

type Analyzer struct {
	files []*os.File
}

func NewAnalyzer(source string) (*Analyzer, error) {
	sourcePath := "testdata/" + source + "/"
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(sourcePath)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			filePaths = append(filePaths, sourcePath+file.Name())
		}
	} else {
		filePaths = append(filePaths, sourcePath+fileInfo.Name())
	}

	analyzer := new(Analyzer)

	for _, filePath := range filePaths {
		jackFile, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		analyzer.files = append(analyzer.files, jackFile)
	}

	return analyzer, nil
}

func (a Analyzer) Files() []*os.File {
	return a.files
}

func (a *Analyzer) Close() {
	for _, file := range a.files {
		file.Close()
	}
}
