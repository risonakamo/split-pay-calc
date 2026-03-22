package splitpay

import (
	"testing"

	"github.com/k0kubun/pp/v3"
)

func Test_readCsv(t *testing.T) {
    result:=readPayCsv("C:/Users/ktkm/Desktop/split-pay-calc/data/data.tsv")

    result2:=csvDataToSplitPayItems(result)

    pp.Println(result2)
}