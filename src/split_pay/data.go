package splitpay

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

// extracted lines from csv
type CsvRow struct {
    Order int `csv:"order"`
    Description string `csv:"description"`

    ChargedCost float64 `csv:"charged amount"`

    OrigPayer string `csv:"payer"`
    SharedPayers string `csv:"shared payers"`
}

// read the split pay csv file
func readPayCsv(filename string) []CsvRow {
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

    return result
}

// normalise payer name
func normaliseName(payer string) string {
    return strings.ToLower(strings.TrimSpace(payer))
}

// convert csv data into split pay items
func csvDataToSplitPayItems(data []CsvRow) []PaymentItem {
    var result []PaymentItem

    var row CsvRow
    for _, row = range data {
        if row.ChargedCost<=0 ||
        len(row.Description)==0 ||
        len(row.OrigPayer)==0 ||
        len(row.SharedPayers)==0 {
            fmt.Println("bad csv data - skipping")
            continue
        }

        var origPayer string = normaliseName(row.OrigPayer)

        // get the split players list
        var splitPayers []string
        var splitPlayersPreClean []string = strings.Split(row.SharedPayers, ",")

        var splitPayer string
        for _, splitPayer = range splitPlayersPreClean {
            var fixedSplitPayer string = normaliseName(splitPayer)

            // ensure original payer is not included
            if fixedSplitPayer == origPayer {
                fmt.Println("split payer was the original payer... strange")
                continue
            }

            splitPayers = append(splitPayers, fixedSplitPayer)
        }

        var item PaymentItem = PaymentItem{
            ItemName:              strings.TrimSpace(row.Description),
            OriginalPayer:         origPayer,
            OriginalPaymentAmount: row.ChargedCost,
            SplitPayers:           splitPayers,
        }

        result = append(result, item)
    }

    return result
}