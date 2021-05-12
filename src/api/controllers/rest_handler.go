package controllers

import (
	"Paxos/src/api/models"
	"Paxos/src/api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type RestHandler struct {
	service services.ServiceProviderInterface
}

func NewRestHandler() *RestHandler {
	return &RestHandler{
		service: services.NewServiceProvider(),
	}
}

func (controller RestHandler) GetAccountBalance(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	balance := controller.service.GetAccountBalance()

	c.JSON(http.StatusOK, balance)
	return
}

func (controller RestHandler) GetTransactionHistory(c *gin.Context) {
	c.Header("Content-type", "application/json")
	history:= controller.service.GetHistory()

	c.JSON(http.StatusOK, history)
	return
}

func (controller RestHandler) GetTransaction(c *gin.Context) {
	c.Header("Content-type", "application/json")
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction ID"})
		return
	}

	transaction, err := controller.service.GetTransaction(id)

	if err != nil {
		if err.Error() == "transaction not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error getting history from storage"})
		return
	}

	c.JSON(http.StatusOK, transaction)
	return
}

func (controller RestHandler) ExecTransaction(c *gin.Context) {
	c.Header("Content-type", "application/json")
	var body *models.TransactionBody
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect body format. Cannot bind input JSON" + err.Error()})
		return
	}

	transaction := controller.service.ExecuteTransaction(body)

	if transaction != nil {
		if transaction.Error() == "invalid transaction: not enough money" {
			c.JSON(http.StatusConflict, gin.H{"message": "invalid transaction: not enough money"})
			return
		}
		if transaction.Error() == "incorrect transaction type" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect transaction type"})
			return
		}
		log.Fatal("There was an unknown server error executing transaction: ", transaction.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error when executing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Succesfully executed transaction"})
	return
}
