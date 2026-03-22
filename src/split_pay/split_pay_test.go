package splitpay

import (
	"testing"

	"github.com/k0kubun/pp/v3"
)

func Test_calcSplit(t *testing.T) {
    var data []PaymentItem = []PaymentItem{
        {
            ItemName: "Dinner",
            OriginalPayer: "Alice",
            OriginalPaymentAmount: 60.0,
            SplitPayers: []string{"Bob", "Charlie"},
        },
        {
            ItemName: "Uber",
            OriginalPayer: "Bob",
            OriginalPaymentAmount: 30.0,
            SplitPayers: []string{"Alice"},
        },
        {
            ItemName: "Snacks",
            OriginalPayer: "Charlie",
            OriginalPaymentAmount: 15.0,
            SplitPayers: []string{"Alice", "Bob"},
        },
        {
            ItemName: "Lunch",
            OriginalPayer: "Dave",
            OriginalPaymentAmount: 40.0,
            SplitPayers: []string{"Alice", "Bob"},
        },
        {
            ItemName: "Gas",
            OriginalPayer: "Dave",
            OriginalPaymentAmount: 50.0,
            SplitPayers: []string{"Charlie"},
        },
        {
            ItemName: "Groceries",
            OriginalPayer: "Dave",
            OriginalPaymentAmount: 90.0,
            SplitPayers: []string{"Alice", "Charlie"},
        },
    }

    result := calculateSplitPayments(data)

    pp.Println(result)
}

func Test_calcSplit2(t *testing.T) {
    var data []PaymentItem = []PaymentItem{
        // p2 and p3 owe 30
        {
            ItemName: "Dinner",
            OriginalPayer: "p1",
            OriginalPaymentAmount: 90,
            SplitPayers: []string{"p2", "p3"},
        },
        // p2 and p3 owe 10
        {
            ItemName: "Dinner2",
            OriginalPayer: "p1",
            OriginalPaymentAmount: 30,
            SplitPayers: []string{"p2", "p3"},
        },
    }

    result := calculateSplitPayments(data)

    pp.Println(result)
}