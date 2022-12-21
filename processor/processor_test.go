package processor

import (
	"fmt"
	"testing"
)

func Test_GetAlphanumericLength(t *testing.T) {
	testString := "abc123&$%"
	expected := 6

	if received := GetAlphanumericLength(testString); received != expected {
		t.Errorf("expected: %d, received: %d", expected, received)
	}
}

func Test_ProcessReceipt(t *testing.T) {
	rp := NewReceiptProcessor()

	tests := []struct {
		receipt Receipt
		want    int
		err     error
	}{
		{
			receipt: Receipt{
				Retailer:     "Walmart",
				PurchaseDate: "2022-12-20",
				PurchaseTime: "20:13",
				Items: []Item{
					{
						ShortDescription: "Coca-Cola Soda Soft Drink, 20 fl oz",
						Price:            "2.08",
					},
					{
						ShortDescription: "Takis Rolls Fuego Tortilla Chips 20 oz",
						Price:            "4.98",
					},
				},
				Total: "7.06",
			},
			want: 12,
		},
		{
			receipt: Receipt{
				Retailer:     "Costco",
				PurchaseDate: "2019-11-13",
				PurchaseTime: "14:01",
				Items: []Item{
					{
						ShortDescription: "   Ribeye Steak   ",
						Price:            "25.00",
					},
				},
				Total: "25.00",
			},
			want: 102,
		},
		{
			receipt: Receipt{
				Retailer:     "Costco",
				PurchaseDate: "2019-11-1300",
				PurchaseTime: "14:01",
				Items: []Item{
					{
						ShortDescription: "   Ribeye Steak   ",
						Price:            "25.00",
					},
				},
				Total: "25.00",
			},
			want: 0,
			err:  fmt.Errorf("purchase date %s is not a valid date", "2019-11-1300"),
		},
		{
			receipt: Receipt{
				Retailer:     "Costco",
				PurchaseDate: "2019-11-13",
				PurchaseTime: "14:01",
				Items: []Item{
					{
						ShortDescription: "   Ribeye Steak   ",
						Price:            "25.00",
					},
				},
				Total: "25.0F",
			},
			want: 0,
			err:  fmt.Errorf("receipt total %s is not a valid dollar amount", "25.0F"),
		},
		{
			receipt: Receipt{
				Retailer:     "Costco",
				PurchaseDate: "2019-11-13",
				PurchaseTime: "14:01",
				Items: []Item{
					{
						ShortDescription: "   Ribeye Steak   ",
						Price:            "25.0F",
					},
				},
				Total: "25.00",
			},
			want: 0,
			err:  fmt.Errorf("item price %s is not a valid dollar amount", "25.0F"),
		},
		{
			receipt: Receipt{
				Retailer:     "Costco",
				PurchaseDate: "2019-11-13",
				PurchaseTime: "1F:01",
				Items: []Item{
					{
						ShortDescription: "   Ribeye Steak   ",
						Price:            "25.00",
					},
				},
				Total: "25.00",
			},
			want: 0,
			err:  fmt.Errorf("purchase time %s is not a valid 24 hour time", "1F:01"),
		},
		{
			receipt: Receipt{
				Retailer:     "Costco",
				PurchaseDate: "2019-11-13",
				PurchaseTime: "10:F1",
				Items: []Item{
					{
						ShortDescription: "   Ribeye Steak   ",
						Price:            "25.00",
					},
				},
				Total: "25.00",
			},
			want: 0,
			err:  fmt.Errorf("purchase time %s is not a valid 24 hour time", "10:F1"),
		},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("Process Receipt Test #%d", i+1)
		t.Run(testname, func(t *testing.T) {
			id, err := rp.ProcessReceipt(tt.receipt)
			if err != tt.err {
				if err == nil || tt.err == nil || err.Error() != tt.err.Error() {
					t.Errorf("expected error: %v, received error: %v", tt.err, err)
					return
				}
			}

			if err != nil {
				return
			}

			points, err := rp.GetReceipt(id)
			if err != nil {
				t.Errorf("unexpected error received: %v", err)
				return
			}

			if points != tt.want {
				t.Errorf("received: %d, want: %d", points, tt.want)
			}
		})
	}
}

func Test_GetReceipt_Present(t *testing.T) {
	rp := NewReceiptProcessor()

	realID, err := rp.ProcessReceipt(Receipt{
		Retailer:     "Walmart",
		PurchaseDate: "2022-12-20",
		PurchaseTime: "20:13",
		Items: []Item{
			{
				ShortDescription: "Coca-Cola Soda Soft Drink, 20 fl oz",
				Price:            "2.08",
			},
			{
				ShortDescription: "Takis Rolls Fuego Tortilla Chips 20 oz",
				Price:            "4.98",
			},
		},
		Total: "7.06",
	})

	if err != nil {
		t.Fatalf("Unexpected test error received: %v", err)
	}

	_, err = rp.GetReceipt(realID)

	if err != nil {
		t.Errorf("expected: %v, received: %v", nil, err)
	}
}

func Test_GetReceipt_Not_Present(t *testing.T) {
	rp := NewReceiptProcessor()

	fakeID := "doesnotexist"
	expectedError := fmt.Errorf("no receipt with ID %s found", fakeID)

	_, err := rp.GetReceipt(fakeID)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("expected error: %v, received: %v", expectedError, err)
	}
}
