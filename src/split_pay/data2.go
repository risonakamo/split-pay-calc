package splitpay

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

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
func ReadPayCsv2(filename string) []CsvRow2 {
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

// convert csv items v2 to payment items
func CsvRow2ToPayItem(rows []CsvRow2) []PaymentItem {
    var result []PaymentItem

    var row CsvRow2
    for _,row = range rows {
        var origPayer string=normaliseName(row.OrigPayer)

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
            OriginalPaymentAmount: row.AmountUsd,
            SplitPayers:           splitPayers,
            DirectPay: row.DirectPay=="yes",
        }

        result = append(result, item)
    }

    return result
}