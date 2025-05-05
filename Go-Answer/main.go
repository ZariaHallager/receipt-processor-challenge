package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var storage = make(map[string]Receipt)
var nextID = 1

func processReceipts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var receipt Receipt
	json.NewDecoder(r.Body).Decode(&receipt)

	receipt.ID = strconv.Itoa(nextID)
	storage[receipt.ID] = receipt
	nextID++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func getPointsForReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/receipts/points/"):]
	receipt := storage[id]
	if receipt.ID == "" {
		http.NotFound(w, r)
		return
	}

	points := calculateAllPoints(&receipt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": points})
}

func getAllReceipts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var allReceipts []Receipt
	for _, receipt := range storage {
		allReceipts = append(allReceipts, receipt)
	}
	json.NewEncoder(w).Encode(allReceipts)
}

func calculateRoundDollarPoints(receipt *Receipt) int {
	dollarPoints := 0
	if receipt.Total[len(receipt.Total)-2:] == "00" {
		dollarPoints += 50
	}
	return dollarPoints
}

func getPoints(receipt *Receipt) int {
	points := 0
	retailer := receipt.Retailer
	for _, char := range retailer {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			points++
		}
	}
	return points
}

func checkIfMultipleOfFive(receipt *Receipt) int {
	points := 0
	totalFloat, _ := strconv.ParseFloat(receipt.Total, 64)
	if int(totalFloat*100)%25 == 0 {
		points += 25
	}
	return points
}

func getPointsFromItems(receipt *Receipt) int {
	points := 0
	for i := 1; i < len(receipt.Items); i += 2 {
		if receipt.Items[i].Price != "" {
			points += 5
		}
	}
	return points
}

func pointsForItemLength(receipt *Receipt) int {
	itemLengthPoints := 0
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			itemLengthPoints += int(priceFloat * 0.2)
		}
	}
	return itemLengthPoints
}

func pointsForOddPurchaseDate(receipt *Receipt) int {
	oddPurchaseDatePoints := 0
	date := receipt.PurchaseDate
	dayStr := date[len(date)-2:]
	day, _ := strconv.Atoi(dayStr)
	if day%2 != 0 {
		oddPurchaseDatePoints += 6
	}
	return oddPurchaseDatePoints
}

func pointsForPurchaseTime(receipt *Receipt) int {
	purchaseTimePoints := 0
	timeParts := strings.Split(receipt.PurchaseTime, ":")
	hours, _ := strconv.Atoi(timeParts[0])
	minutes, _ := strconv.Atoi(timeParts[1])
	if float64(hours) == 14.0 && (float64(minutes) > 0.0 && float64(minutes) < 60.0) {
		purchaseTimePoints += 10
	} else if float64(hours) == 15.0 {
		purchaseTimePoints += 10
	}
	return purchaseTimePoints
}

func calculateAllPoints(receipt *Receipt) int {
	totalPoints := 0
	totalPoints += getPoints(receipt)
	totalPoints += checkIfMultipleOfFive(receipt)
	totalPoints += getPointsFromItems(receipt)
	totalPoints += pointsForItemLength(receipt)
	totalPoints += pointsForOddPurchaseDate(receipt)
	totalPoints += pointsForPurchaseTime(receipt)
	totalPoints += calculateRoundDollarPoints(receipt)
	return totalPoints
}

func main() {
	http.HandleFunc("/receipts/process", processReceipts)
	http.HandleFunc("/receipts/points/", getPointsForReceipt)
	http.HandleFunc("/", getAllReceipts)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
