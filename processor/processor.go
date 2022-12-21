package processor

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Item holds the information for a single item in a receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Receipt holds the the receipt data to be processed
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// ReceiptProcessor processes new receipts and stores the points in the receipts map
type ReceiptProcessor struct {
	receipts map[string]int
}

// NewReceiptProcessor returns a new empty receipt processor
func NewReceiptProcessor() ReceiptProcessor {
	r := make(map[string]int)

	return ReceiptProcessor{
		receipts: r,
	}
}

// GetAlphanumericLength counts the number of alphanumeric characters in a string and returns the result
func GetAlphanumericLength(s string) int {
	len := 0
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			continue
		}
		len++
	}

	return len
}

// ProcessReceipt will accept a Receipt argument, process the receipt, calculate and store the points.
// If successful, the generated receipt ID is returned with no error. Else, an empty string and the
// corresponding error are returned.
func (rp *ReceiptProcessor) ProcessReceipt(receipt Receipt) (string, error) {
	points := 0
	// One point for every alphanumeric character in the retailer name.
	points += GetAlphanumericLength(receipt.Retailer)
	// 50 points if the total is a round dollar amount with no cents.
	if receipt.Total[len(receipt.Total)-2:] == "00" {
		points += 50
	}
	// 25 points if the total is a multiple of 0.25.
	f, err := strconv.ParseFloat(receipt.Total, 32)
	if err != nil {
		return "", fmt.Errorf("receipt total %v is not a valid dollar amount", receipt.Total)
	}
	if int(math.Ceil(f*100))%25 == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += 5 * (len(receipt.Items) / 2)

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2
	// and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		if len(strings.Trim(item.ShortDescription, " "))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 32)
			if err != nil {
				return "", fmt.Errorf("item price %v is not a valid dollar amount", item.Price)
			}

			points += int(math.Ceil(price * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd.
	date, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return "", fmt.Errorf("purchase date %s is not a valid date", receipt.PurchaseDate)
	}
	if date.Day()%2 == 1 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	prefix, err := strconv.Atoi(receipt.PurchaseTime[0:2])
	if err != nil {
		return "", fmt.Errorf("purchase time %s is not a valid 24 hour time", receipt.PurchaseTime)
	}
	suffix, err := strconv.Atoi(receipt.PurchaseTime[len(receipt.PurchaseTime)-2:])
	if err != nil {
		return "", fmt.Errorf("purchase time %s is not a valid 24 hour time", receipt.PurchaseTime)
	}

	if (prefix == 14 && suffix > 0) || (prefix > 14 && prefix < 16) {
		points += 10
	}

	id := uuid.New().String()
	rp.receipts[id] = points
	return id, nil
}

// GetReceipt will check if the requested receipt exists. If it does, it will return the points;
// if not, it will return an error
func (rp ReceiptProcessor) GetReceipt(receiptID string) (int, error) {
	points, ok := rp.receipts[receiptID]

	if !ok {
		return 0, fmt.Errorf("no receipt with ID %s found", receiptID)
	}
	return points, nil
}
