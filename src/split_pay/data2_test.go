package splitpay

import (
	"testing"

	"github.com/k0kubun/pp/v3"
)

func Test_dataRead(t *testing.T) {
    result:=readPayCsv2("C:/Users/ktkm2/Desktop/newprojs/split-pay-calc/data/data2.tsv")

    pp.Println(result)
}