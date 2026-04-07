package splitpay

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/gocarina/gocsv"
)

// csv row from sheet format
type CsvRow2 struct {
	Order       int    `csv:"order"`
	Description string `csv:"item"`

	AmountYen string  `csv:"amount (yen)"`
	AmountUsd float64 `csv:"amount (usd)"`

	OrigPayer    string `csv:"original payer"`
	SharedPayers string `csv:"shared payers"`
	DirectPay    string `csv:"direct pay"`
}

// read pay csv file v2
func readPayCsv2(filename string) []CsvRow2 {
	var e error
    var file *os.File
	file,e=os.Open(filename)
    if e!=nil {
        panic(e)
    }
    defer file.Close()

    var result []CsvRow2

    gocsv.SetCSVReader(func(ioReader io.Reader) gocsv.CSVReader {
        var newCsvReader *csv.Reader=csv.NewReader(ioReader)
        newCsvReader.Comma='\t'
        newCsvReader.FieldsPerRecord=-1
        return newCsvReader
    })

    e=gocsv.UnmarshalFile(file,&result)
    if e!=nil {
        panic(e)
    }

    return result
}