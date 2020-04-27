package test

import (
	"github.com/stretchr/testify/assert"
	"paxos/src/api/models"
	"paxos/src/api/services"
	"testing"
)

func TestGetAccountBalance(t *testing.T) {
	serviceProvider := services.NewServiceProvider()

	balance := serviceProvider.GetAccountBalance()

	assert.NotNil(t, balance)
}

func TestGetHistory(t *testing.T) {
	serviceProvider := services.NewServiceProvider()

	history, err := serviceProvider.GetHistory()

	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, history)
}

func TestGetTransaction(t *testing.T) {
	serviceProvider := services.NewServiceProvider()

	transaction, err := serviceProvider.GetTransaction("1")

	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, transaction)
}

func TestGetTransaction_TransactionNotFound(t *testing.T) {
	serviceProvider := services.NewServiceProvider()

	transaction, err := serviceProvider.GetTransaction("209")

	assert.Nil(t, transaction)
	assert.Error(t, err)
	assert.EqualError(t, err, "transaction not found")
}

func TestExecuteTransaction_Credit(t *testing.T) {
	serviceProvider := services.NewServiceProvider()

	transaction := &models.TransactionBody{
		Type: "credit",
		Ammount: 10,

	}

	err := serviceProvider.ExecuteTransaction(transaction)

	assert.NoError(t, err)
	assert.Nil(t, err)
}

func TestExecuteTransaction_Debit(t *testing.T) {
	serviceProvider := services.NewServiceProvider()

	transaction := &models.TransactionBody{
		Type: "debit",
		Ammount: 10,

	}

	err := serviceProvider.ExecuteTransaction(transaction)

	assert.NoError(t, err)
	assert.Nil(t, err)
}

func TestExecuteTransaction_Debit_InvalidAmmount_NotEnoughMoney(t *testing.T) {
	serviceProvider := services.NewServiceProvider()

	transaction := &models.TransactionBody{
		Type: "debit",
		Ammount: 10100.00,

	}

	err := serviceProvider.ExecuteTransaction(transaction)

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid transaction: not enough money")
}

func TestExecuteTransaction_InvalidTransactionType(t *testing.T) {
	serviceProvider := services.NewServiceProvider()

	transaction := &models.TransactionBody{
		Type: "invalid",
		Ammount: 20,

	}

	err := serviceProvider.ExecuteTransaction(transaction)

	assert.Error(t, err)
	assert.EqualError(t, err, "incorrect transaction type")
}
