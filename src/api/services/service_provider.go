package services

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"paxos/src/api/database"
	"paxos/src/api/models"
	"strconv"
	"time"
)

type ServiceProviderInterface interface {
	GetAccountBalance() models.AccountBalance
	GetHistory() ([]models.Transaction, error)
	GetTransaction(id string) (*models.Transaction, error)
	ExecuteTransaction(transaction *models.TransactionBody) error
}

type ServiceProvider struct {
	// IF DB IS USED FOR PERSISTENCE, HERE WE DECLARE THE DB HANDLER WE WOULD USE FOR THE SERVICE ( DATABASE CONFIGURATION WOULD BE STORED ELSEWHERE)
	//db   *sql.DB
	//stmt *sql.Stmt

	// Since no db is used, we emulate db as follows
	dbHandler database.DbInterface
	Balance float64

}

// INSTANCE SERVICE PROVIDER

func NewServiceProvider() ServiceProviderInterface {
	return ServiceProvider{
		dbHandler: database.NewDbHandler(),
		Balance: 0,
	}
}

func (service ServiceProvider) GetAccountBalance() models.AccountBalance {
	response := models.AccountBalance{
		Balance: service.Balance,
	}
	return response
}
func (service ServiceProvider) GetHistory() ([]models.Transaction, error) {
	var response []models.Transaction
	data, err := service.dbHandler.ReadFile()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	errMarshal := json.Unmarshal(data, &response)
	if errMarshal != nil {
		log.Fatal(errMarshal)
		return nil, errMarshal
	}

	return response, nil
}

func (service ServiceProvider) GetTransaction(id string) (*models.Transaction, error) {
	var transactionHistory []models.Transaction

	data, err := service.dbHandler.ReadFile()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	errMarshal := json.Unmarshal(data, &transactionHistory)
	if errMarshal != nil {
		log.Fatal(errMarshal)
		return nil, errMarshal
	}

	for key,element := range transactionHistory {
		if element.ID == id {
			return &transactionHistory[key], nil
		}
	}

	return nil, errors.New("transaction not found")

}

func (service ServiceProvider) ExecuteTransaction(transaction *models.TransactionBody) error {
	switch transaction.Type {
	case "debit":
		result := service.Balance - transaction.Ammount

		if result < 0 {
			return errors.New("invalid transaction: not enough money")
		}

		service.Balance = result

	case "credit":
		result := service.Balance + transaction.Ammount

		service.Balance = result
	default:
		return errors.New("incorrect transaction type")
	}

	dateStamp := time.Now()

	transactionToWrite := models.Transaction{
		// Random INT is not recommended because cannot ensure PK standards. Perfect solution for this would be transactional DB usage
		ID:      strconv.Itoa(rand.Int()),
		Type:    transaction.Type,
		Ammount: transaction.Ammount,
		Date:    dateStamp.String(),
	}

	// Convert transaction to []byte
	data, err := json.Marshal(&transactionToWrite)
	if err != nil {
		log.Fatal(err)
		return  err
	}

	// Write new transaction in file (aux db)
	err = service.dbHandler.WriteFile(data)
	if err != nil {
		log.Fatal(err)
		return  err
	}

	return nil
}
