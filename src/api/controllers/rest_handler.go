package controllers

import (
	"github.com/gin-gonic/gin"
	"paxos/src/api/models"
	"log"
	"net/http"
	"paxos/src/api/services"
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
}

func (controller RestHandler) GetTransactionHistory(c *gin.Context) {
	c.Header("Content-type", "application/json")
	history, err := controller.service.GetHistory()

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error getting history from storage"})
		return
	}



	c.JSON(http.StatusOK, history)
	return
}

func (controller RestHandler) GetTransaction(c *gin.Context) {
	c.Header("Content-type", "application/json")
	id := c.Param("id")

	// I am assuming that ID, since it is a string, could not only be numeric. If assumption is incorrect, would check if id is not a number with following commented code

	//if _, err := strconv.Atoi(v); err == nil {
	//	response should be 400 Bad request
	//}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction ID"})
	}

	transaction, err := controller.service.GetTransaction(id)

	if err != nil {
		if err.Error() == "transaction not found" {
			log.Println("Error : Transaction not found")
			c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		}
		log.Fatal(err)
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
			log.Println("invalid transaction: not enough money")
			c.JSON(http.StatusConflict, gin.H{"message": "invalid transaction: not enough money"})
			return
		}
		if transaction.Error() == "incorrect transaction type" {
			log.Println("incorrect transaction type")
			c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect transaction type"})
			return
		}
		log.Fatal("There was an unknown server error executing transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "There was an error when executing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Succesfully executed transaction"})
	return
}
