// package implementing split pay algo and data types

package splitpay

import (
	"split-pay-calc/src/utils"

	mapset "github.com/deckarep/golang-set/v2"
)

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

    // shared payers pay directly to the original payer - original payer does
    // not split with the shared payers
    DirectPay bool
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
    Itemisation []ItemisedPayItem `yaml:"items"`
}

// shorter form of pay item for itemisation list
type ItemisedPayItem struct {
    ItemName string `yaml:"description"`

    // can be negative. represents a payment detracting from the sub-payer
    // owed amount
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

        // remove original payer if direct pay
        if item.DirectPay {
            totalPayers-=1
        }

        // the split payment of all payers. all split payers must pay this amount to
        // the original payer (who has already paid this split payment)
        var splitPayment float64=item.OriginalPaymentAmount/float64(totalPayers)

        var subPayer string
        for _,subPayer = range item.SplitPayers {
            // dict initialisations
            var payerInResult bool
            _,payerInResult=payersResultDict[subPayer]

            if !payerInResult {
                payersResultDict[subPayer]=&SplitPayResult{
                    Payer: subPayer,
                    Payments: PaymentsDict{},
                }
            }

            var origPayerInSubpayerPayments bool
            _,origPayerInSubpayerPayments=payersResultDict[subPayer].
                Payments[item.OriginalPayer]

            if !origPayerInSubpayerPayments {
                payersResultDict[subPayer].Payments[item.OriginalPayer]=&Payment{
                    ToPerson: item.OriginalPayer,
                    Amount: 0,
                }
            }

            // ensuring the reverse direction exists
            var in bool
            _,in=payersResultDict[item.OriginalPayer]

            if !in {
                payersResultDict[item.OriginalPayer]=&SplitPayResult{
                    Payer: subPayer,
                    Payments: PaymentsDict{},
                }
            }

            _,in=payersResultDict[item.OriginalPayer].Payments[subPayer]

            if !in {
                payersResultDict[item.OriginalPayer].Payments[subPayer]=&Payment{
                    ToPerson: item.OriginalPayer,
                    Amount: 0,
                }
            }



            // filling out payer->subpayer
            payersResultDict[subPayer].Payments[item.OriginalPayer].Amount+=splitPayment
            payersResultDict[subPayer].Total+=splitPayment
            payersResultDict[subPayer].Payments[item.OriginalPayer].Itemisation=append(
                payersResultDict[subPayer].Payments[item.OriginalPayer].Itemisation,
                ItemisedPayItem{
                    ItemName: item.ItemName,
                    Amount: splitPayment,
                },
            )

            // filling out reverse (just the itemisation)
            payersResultDict[item.OriginalPayer].Payments[subPayer].Itemisation=append(
                payersResultDict[item.OriginalPayer].Payments[subPayer].Itemisation,
                ItemisedPayItem{
                    ItemName: item.ItemName,
                    Amount: -splitPayment,
                },
            )
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

// balance split payer dict to come up with simplest payment plan
func BalanceSplitPay(splitPay SplitPayersDict) SplitPayersDict {
    var allNames []string=getAllNames(splitPay)

    var name1 string
    var name2 string

    var in bool

    for _,name1 = range allNames {
        for _,name2 = range allNames {
            // extracting amounts for name1 and name2
            var amount1 float64=0
            var amount2 float64=0

            _,in=splitPay[name1]

            if in {
                _,in=splitPay[name1].Payments[name2]

                if in {
                    amount1=splitPay[name1].Payments[name2].Amount
                }
            }

            _,in=splitPay[name2]

            if in {
                _,in=splitPay[name2].Payments[name1]

                if in {
                    amount2=splitPay[name2].Payments[name1].Amount
                }
            }

            // if either amount is 0, do nothing. nothing to change on either side
            if amount1==0 || amount2==0 {
                continue
            }

            // amount 1 is less than 2. set amount1 to zero and deduct from amount 2
            if amount1<amount2 {
                splitPay[name1].Payments[name2].Amount=0
                splitPay[name1].Total-=amount1

                splitPay[name2].Payments[name1].Amount-=amount1
                splitPay[name2].Total-=amount1

            // otherwise, do the reverse
            } else {
                splitPay[name2].Payments[name1].Amount=0
                splitPay[name2].Total-=amount2

                splitPay[name1].Payments[name2].Amount-=amount2
                splitPay[name1].Total-=amount2
            }
        }
    }

    return splitPay
}

// extract all possible names from split pay dict
func getAllNames(splitPay SplitPayersDict) []string {
    var result mapset.Set[string]=mapset.NewSet[string]()

    var name string
    for name,_ = range splitPay {
        result.Add(name)

        var name2 string
        for name2,_ = range splitPay[name].Payments {
            result.Add(name2)
        }
    }

    return result.ToSlice()
}