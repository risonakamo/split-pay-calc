// package implementing split pay algo and data types

package splitpay

import "split-pay-calc/src/utils"

// split pay result keyed by the payer
type SplitPayersDict map[string]*SplitPayResult

// payments keyed by the receiver of the payment
type PaymentsDict map[string]*Payment

// dict with total payments by original payers
type TotalsDict map[string]float64

// a payment that occurred that needs to be split
type PaymentItem struct {
    ItemName string
    OriginalPayer string
    OriginalPaymentAmount float64

    // list of other payers which this payment should be split with.
    // this should not include the original payer.
    SplitPayers []string
}

// payment result for a single person. contains the payments that the person needs
// to make, and to whom
type SplitPayResult struct {
    Payer string `yaml:"payer"`
    Payments PaymentsDict `yaml:"payments"`
    Total float64 `yaml:"total"`
}

// a payment to be made to a certain person
type Payment struct {
    ToPerson string `yaml:"toPerson"`
    Amount float64 `yaml:"amount"`
}

// container for split pays results
type SplitPayResultTop struct {
    SplitPays SplitPayersDict `yaml:"split"`
    Totals TotalsDict `yaml:"totalOriginalPayments"`
}

// given payment items, calculate the payment splits for all payers involved in the items
func CalculateSplitPayments(items []PaymentItem) SplitPayResultTop {
    var payersResultDict SplitPayersDict=SplitPayersDict{}
    var totals TotalsDict=TotalsDict{}

    var item PaymentItem
    for _,item = range items {
        // adding to totals dict
        var inTotals bool
        _,inTotals=totals[item.OriginalPayer]

        if !inTotals {
            totals[item.OriginalPayer]=0
        }

        totals[item.OriginalPayer]+=item.OriginalPaymentAmount

        // total payers of the item, which includes the original payer
        var totalPayers int=len(item.SplitPayers)+1

        // the split payment of all payers. all split payers must pay this amount to
        // the original payer (who has already paid this split payment)
        var splitPayment float64=item.OriginalPaymentAmount/float64(totalPayers)

        var subPayer string
        for _,subPayer = range item.SplitPayers {
            var payerInResult bool
            _,payerInResult=payersResultDict[subPayer]

            if !payerInResult {
                payersResultDict[subPayer]=&SplitPayResult{
                    Payer: subPayer,
                    Payments: PaymentsDict{},
                }
            }

            var origPayerInSubpayerPayments bool
            _,origPayerInSubpayerPayments=payersResultDict[subPayer].Payments[item.OriginalPayer]

            if !origPayerInSubpayerPayments {
                payersResultDict[subPayer].Payments[item.OriginalPayer]=&Payment{
                    ToPerson: item.OriginalPayer,
                    Amount: 0,
                }
            }

            payersResultDict[subPayer].Payments[item.OriginalPayer].Amount+=splitPayment
            payersResultDict[subPayer].Total+=splitPayment
        }
    }

    return SplitPayResultTop{
        SplitPays: payersResultDict,
        Totals: totals,
    }
}

// round all relevant float vals in result
func RoundSplitResult(result SplitPayResultTop) SplitPayResultTop {
    // round split pays
    var splitPay *SplitPayResult
    for _, splitPay = range result.SplitPays {
        // round total
        splitPay.Total = utils.ToFixed(splitPay.Total, 2)

        // round each payment
        var payment *Payment
        for _, payment = range splitPay.Payments {
            payment.Amount = utils.ToFixed(payment.Amount, 2)
        }
    }

    // round totals dict
    var payer string
    var payVal float64
    for payer,payVal = range result.Totals {
        result.Totals[payer] = utils.ToFixed(payVal, 2)
    }

    return result
}