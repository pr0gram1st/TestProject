package main

import (
	"database/sql"
	_ "fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Account struct {
	ID      int
	Name    string
	Balance float64
}

type Transaction struct {
	ID         int
	Value      float64
	AccountID  int
	Group      string
	Account2ID int
	Date       time.Time
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=0000 dbname=testDB sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.POST("/account", createAccount)
	router.GET("/account/:id", getAccount)
	router.GET("/accounts", getAllAccounts)
	router.POST("/transaction", createTransaction)
	router.GET("/transaction/:id", getTransactions)

	router.Run(":8080")
}

func createAccount(c *gin.Context) {
	var account Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `INSERT INTO accounts (name, balance) VALUES ($1, $2) RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, account.Name, account.Balance).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getAccount(c *gin.Context) {
	id := c.Param("id")

	var account Account
	sqlStatement := `SELECT id, name, balance FROM accounts WHERE id=$1;`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&account.ID, &account.Name, &account.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func getAllAccounts(c *gin.Context) {
	var accounts []Account
	sqlStatement := `SELECT id, name, balance FROM accounts;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		err = rows.Scan(&account.ID, &account.Name, &account.Balance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, accounts)
}

func createTransaction(c *gin.Context) {
	var transaction Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if transaction.Group == "transfer" && transaction.Account2ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account2ID is required for transfers"})
		return
	}

	sqlStatement := `INSERT INTO transactions (value, accountid, ggroup, account2id, datet) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, transaction.Value, transaction.AccountID, transaction.Group, transaction.Account2ID, transaction.Date).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getTransactions(c *gin.Context) {
	id := c.Param("id")

	var transactions []Transaction
	sqlStatement := `SELECT id, value, accountid, ggroup, account2id, datet FROM transactions WHERE accountid=$1;`
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		err = rows.Scan(&transaction.ID, &transaction.Value, &transaction.AccountID, &transaction.Group, &transaction.Account2ID, &transaction.Date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}
