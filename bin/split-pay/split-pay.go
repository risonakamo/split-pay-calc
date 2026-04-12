package main

import (
	"path/filepath"
	splitpay "split-pay-calc/src/split_pay"
	"split-pay-calc/src/utils"
)

func main() {
    var here string=utils.GetHereDirExe()

    var dataPath string=filepath.Join(here,"../../data/data3.tsv")

    var payItems []splitpay.PaymentItem=splitpay.CsvRow2ToPayItem(
        splitpay.ReadPayCsv2(dataPath),
    )

    var result splitpay.SplitPayResultTop=splitpay.CalculateSplitPayments(payItems)

    result.SplitPays=splitpay.BalanceSplitPay(result.SplitPays)

    result=splitpay.RoundSplitResult(result)

    utils.WriteYaml(
        filepath.Join(here,"out.yml"),
        result,
    )
}