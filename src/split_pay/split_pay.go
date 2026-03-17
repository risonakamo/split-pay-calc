// package implementing split pay algo and data types

package splitpay

// split pay result keyed by the payer
type SplitPayersDict map[string]SplitPayResult

// a payment that occurred that needs to be split
type PaymentItem struct {
    ItemName string
    OriginalPayer string
    OriginalPaymentAmount float32

    // list of other payers which this payment should be split with.
    // this should not include the original payer.
    SplitPayers []string
}

// payment result for a single person. contains the payments that the person needs
// to make, and to whom
type SplitPayResult struct {
    Payer string
    Payments []Payment
}

// a payment to be made to a certain person
type Payment struct {
    ToPerson string
    Amount float32
}

// given payment items, calculate the payment splits for all payers involved in the items
func calculateSplitPayments(items []PaymentItem) []SplitPayResult {
    var item PaymentItem
    for _,item = range items {
        // total payers of the item, which includes the original payer
        var totalPayers int=len(item.SplitPayers)+1
        // the split payment of all payers. all split payers must pay this amount to
        // the original payer (who has already paid this split payment)
        var splitPayment float=item.OriginalPaymentAmount/float32(totalPayers)


    }
}