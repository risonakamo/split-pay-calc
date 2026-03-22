package utils

import (
	"math"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// give folder location of the exe that calls this func
func GetHereDirExe() string {
    var exePath string
    var e error
    exePath,e=os.Executable()

    if e!=nil {
        panic(e)
    }

    return filepath.Dir(exePath)
}

// overwrite target yml file with a new file
func WriteYaml(filename string,data any) error {
	var wfile *os.File
	var e error
	wfile,e=os.Create(filename)

	if e!=nil {
		panic(e)
	}

	defer wfile.Close()

	var ymldata []byte
	ymldata,e=yaml.Marshal(data)

	if e!=nil {
		panic(e)
	}

	wfile.Write(ymldata)
	return nil
}

// round to int
func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

// round to decimal point
func ToFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}