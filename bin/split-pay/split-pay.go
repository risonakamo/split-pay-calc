package main

import (
	"path/filepath"
	splitpay "split-pay-calc/src/split_pay"
	"split-pay-calc/src/utils"
)

func main() {
    var here string=utils.GetHereDirExe()

    var dataPath string=filepath.Join(here,"../../data/data.tsv")

    var payItems []splitpay.PaymentItem=splitpay.ReadPayCsvToItems(dataPath)

    var result splitpay.SplitPayResultTop=splitpay.CalculateSplitPayments(payItems)

    utils.WriteYaml(
        filepath.Join(here,"out.yml"),
        result,
    )
}