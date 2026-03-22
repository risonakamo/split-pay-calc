package splitpay

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/k0kubun/pp/v3"
)

// extracted lines from csv
type CsvRow struct {
    Order int `csv:"order"`
    Description string `csv:"description"`

    ChargedCost float32 `csv:"charged amount"`

    OrigPayer string `csv:"payer"`
    SharedPayers string `csv:"shared payers"`
}

// read the split pay csv file
func readPayCsv(filename string) {
	var e error
    var file *os.File
	file,e=os.Open(filename)
    if e!=nil {
        panic(e)
    }
    defer file.Close()

    var result []CsvRow

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

    pp.Println(result)
}