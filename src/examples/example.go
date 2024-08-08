package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"touch-coding-challenge/src/models"
)

const baseURL = "http://localhost:8080"

// Helper function to send HTTP requests and print responses
func sendRequest(method, endpoint string, body interface{}) (*http.Response, string, error) {
	var reqBody io.Reader
	url := baseURL + endpoint
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, "", err
		}
		reqBody = strings.NewReader(string(jsonBody))
	} else {
		reqBody = nil
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, "", err
	}

	return resp, string(bodyBytes), nil
}

func main() {
	// -----------------------------------------------------------
	// Create client 1 with 1 receipt
	client_id := "12345"
	client := models.Client{
		ClientID: client_id,
		Name:     "John Doe",
	}
	_, body, err := sendRequest("POST", "/clients", client)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	fmt.Printf("Create Client Response: %s\n", body)

	// Create a new deposit
	depositID := uuid.NewString()

	var receipts []models.Receipt
	receipt := models.Receipt{
		Amount:      float64(10000),
		Pot:         "Pot A",
		WrapperType: "GIA",
	}
	receipts = append(receipts, receipt)

	deposit := models.Deposit{
		DepositID: depositID,
		ClientID:  client_id,
		Receipts:  receipts,
	}
	_, body, err = sendRequest("POST", "/deposits", deposit)
	if err != nil {
		log.Fatalf("Error creating deposit: %v", err)
	}
	fmt.Printf("Create Deposit Response: %s\n", body)

	// Retrieve the deposit which has a single receipt
	_, body, err = sendRequest("GET", fmt.Sprintf("/deposits/%s", depositID), nil)
	if err != nil {
		log.Fatalf("Error retrieving deposit: %v", err)
	}
	fmt.Printf("Retrieve Deposit Response: %s\n", body)

	// -----------------------------------------------------------
	// Send a second deposit with 3 receipts and a second pot
	depositID = uuid.NewString()

	receipt.Amount = float64(20000)
	receipt.WrapperType = "ISA"
	receipts = append(receipts, receipt)

	receipt.Amount = float64(50000)
	receipt.WrapperType = "SIPP"
	receipts = append(receipts, receipt)

	receipt.Amount = float64(20000)
	receipt.Pot = "Pot B"
	receipt.WrapperType = "GIA"
	receipts = append(receipts, receipt)

	deposit = models.Deposit{
		DepositID: depositID,
		ClientID:  client_id,
		Receipts:  receipts,
	}
	_, body, err = sendRequest("POST", "/deposits", deposit)
	if err != nil {
		log.Fatalf("Error creating deposit: %v", err)
	}
	fmt.Printf("Create Deposit Response: %s\n", body)

	// Retrieve the deposits which has 4 receipts
	_, body, err = sendRequest("GET", fmt.Sprintf("/deposits/%s", depositID), nil)
	if err != nil {
		log.Fatalf("Error retrieving deposit: %v", err)
	}
	fmt.Printf("Retrieve Deposit Response: %s\n", body)

	// -----------------------------------------------------------
	// Create a second client
	client = models.Client{
		ClientID: "67890",
		Name:     "John Doe 2",
	}
	_, body, err = sendRequest("POST", "/clients", client)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	fmt.Printf("Create Client Response: %s\n", body)

	// Try to deposit into a client that does not exist
	deposit.ClientID = "1234567890"
	depositID = uuid.NewString()
	deposit.DepositID = depositID

	_, body, err = sendRequest("POST", "/deposits", deposit)
	if err != nil {
		log.Fatalf("Error creating deposit: %v", err)
	}
	fmt.Printf("Create Deposit Response: %s\n", body)

	// Deposit into a client exists
	deposit.ClientID = "67890"
	_, body, err = sendRequest("POST", "/deposits", deposit)
	if err != nil {
		log.Fatalf("Error creating deposit: %v", err)
	}
	fmt.Printf("Create Deposit Response: %s\n", body)

	// Retrieve the deposits which has 3 receipts
	_, body, err = sendRequest("GET", fmt.Sprintf("/deposits/%s", depositID), nil)
	if err != nil {
		log.Fatalf("Error retrieving deposit: %v", err)
	}
	fmt.Printf("Retrieve Deposit Response: %s\n", body)
}
