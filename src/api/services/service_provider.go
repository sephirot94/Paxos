package services

import (
	"Paxos/src/api/database"
	"Paxos/src/api/models"
	"errors"
	"math/rand"
	"strconv"
	"time"
)

type ServiceProviderInterface interface {
	GetAccountBalance() models.AccountBalance
	GetHistory() []models.Transaction
	GetTransaction(id string) (*models.Transaction, error)
	ExecuteTransaction(transaction *models.TransactionBody) error
}

type ServiceProvider struct {
	// IF DB IS USED FOR PERSISTENCE, HERE WE DECLARE THE DB HANDLER WE WOULD USE FOR THE SERVICE ( DATABASE CONFIGURATION WOULD BE STORED ELSEWHERE)
	//db   *sql.DB
	//stmt *sql.Stmt

	// Since no db is used, we emulate db as follows
	dbHandler database.DbInterface
}

// INSTANCE SERVICE PROVIDER

func NewServiceProvider() ServiceProviderInterface {
	return ServiceProvider{
		dbHandler: database.NewDbHandler(),
	}
}

func (service ServiceProvider) GetAccountBalance() models.AccountBalance {
	data := service.dbHandler.GetData()
	return data.Balance
}

func (service ServiceProvider) GetHistory() []models.Transaction {
	data := service.dbHandler.GetData()

	return data.History
}

func (service ServiceProvider) GetTransaction(id string) (*models.Transaction,error) {

	data := service.dbHandler.GetData()

	for key,element := range data.History {
		if element.ID == id {
			return &data.History[key], nil
		}
	}

	return nil, errors.New("transaction not found")

}

func (service ServiceProvider) ExecuteTransaction(transaction *models.TransactionBody) error {
	data := service.dbHandler.GetData()
	switch transaction.Type {
	case "debit":
		result := data.Balance.Balance - transaction.Ammount

		if result < 0 {
			return errors.New("invalid transaction: not enough money")
		}

	case "credit":

	default:
		return errors.New("incorrect transaction type")
	}

	dateStamp := time.Now()

	newTransaction := models.Transaction{
		// Random INT is not recommended because cannot ensure PK standards. Perfect solution for this would be transactional DB usage
		ID:      strconv.Itoa(rand.Int()),
		Type:    transaction.Type,
		Ammount: transaction.Ammount,
		Date:    dateStamp.String(),
	}

	service.dbHandler.SetData(newTransaction)

	return nil
}
