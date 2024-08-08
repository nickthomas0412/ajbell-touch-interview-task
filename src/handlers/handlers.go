package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"touch-coding-challenge/src/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var db *sql.DB

// SetDatabase sets the DB var so it can be used in the handler package
func SetDatabase(database *sql.DB) {
	db = database
}

// CreateClient creates a client entry in the table if it is not present
func CreateClient(w http.ResponseWriter, r *http.Request) {
	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientID := client.ClientID
	query := `SELECT client_id FROM Clients WHERE client_id = ?`
	err := db.QueryRow(query, clientID).Scan(&client.ClientID)
	if err == sql.ErrNoRows {
		query = `INSERT INTO Clients (client_id, name) VALUES (?, ?)`
		_, err := db.Exec(query, client.ClientID, client.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

// CreateDeposit creates a deposit entry in the table and subsequent entries for all receipts
func CreateDeposit(w http.ResponseWriter, r *http.Request) {
	var deposit models.Deposit
	var amount float64
	if err := json.NewDecoder(r.Body).Decode(&deposit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if client already exists otherwise return error
	clientID := deposit.ClientID
	query := `SELECT client_id FROM Clients WHERE client_id = ?`
	err := db.QueryRow(query, clientID).Scan(&clientID)
	if err == sql.ErrNoRows {
		http.Error(w, "Client does not exist", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add deposit to deposit table
	query = `INSERT INTO Deposits (deposit_id, client_id, nominal) VALUES (?, ?, ?)`
	receipts := deposit.Receipts
	for i := range receipts {
		amount += receipts[i].Amount
	}

	_, err = db.Exec(query, deposit.DepositID, deposit.ClientID, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add receipts to receipt table for the whole deposit
	query = `INSERT INTO Receipts (receipt_id, deposit_id, amount, pot, wrapper_type) VALUES (?, ?, ?, ?, ?)`
	for _, receipt := range receipts {
		receipt.ReceiptID = uuid.NewString()
		receipt.DepositID = deposit.DepositID
		_, err := db.Exec(query, receipt.ReceiptID, receipt.DepositID, receipt.Amount, receipt.Pot, receipt.WrapperType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deposit)
}

// RetrieveDeposit returns all receipts for an existing deposit
func RetrieveDeposit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	depositID := vars["id"]

	var deposit models.Deposit
	query := `SELECT deposit_id, client_id, nominal, created_at FROM Deposits WHERE deposit_id = ?`
	err := db.QueryRow(query, depositID).Scan(&deposit.DepositID, &deposit.ClientID, &deposit.Nominal, &deposit.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Deposit not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	rows, err := db.Query(`SELECT receipt_id, deposit_id, amount, pot, wrapper_type, created_at FROM Receipts WHERE deposit_id = ?`, depositID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var receipts []models.Receipt
	for rows.Next() {
		var receipt models.Receipt
		if err := rows.Scan(&receipt.ReceiptID, &receipt.DepositID, &receipt.Amount, &receipt.Pot, &receipt.WrapperType, &receipt.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		receipts = append(receipts, receipt)
	}
	deposit.Receipts = receipts

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deposit)
}
