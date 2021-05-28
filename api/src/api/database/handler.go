package database

import (
	"Paxos/src/api/models"
	"sync"
)

var Data = models.DB{
	Balance: models.AccountBalance{
		Balance: 100,
	},
	History: []models.Transaction{},
	}

type DbInterface interface {
	GetData() models.DB
	SetData(transaction models.Transaction)
}

type DbHandler struct {
	Mutex *sync.RWMutex
	WaitGroup *sync.WaitGroup
}

func NewDbHandler() DbInterface {
	var m sync.RWMutex
	var wg sync.WaitGroup
	return DbHandler{
		Mutex: &m,
		WaitGroup: &wg,
	}
}

//func (h DbHandler) ReadFile() ([]byte, error) {
//	// Check if file exists
//	_, err := os.Stat("/Users/ijinkus/go/src/paxos/src/api/database/db.json")
//	if err != nil {
//		if os.IsNotExist(err) {
//			log.Fatal("File does not exist.")
//			return nil, err
//		}
//	}
//	// Wait for pending write operations
//	h.WaitGroup.Wait()
//
//	data, err := ioutil.ReadFile("/Users/ijinkus/go/src/paxos/src/api/database/db.json")
//	if err != nil {
//		log.Fatal(err)
//		return nil, err
//	}
//
//	return data, nil
//}
//
//func (h DbHandler) WriteFile(data []byte) error {
//	// Check if file exists
//	_, err := os.Stat("/Users/ijinkus/go/src/paxos/src/api/database/db.json")
//	if err != nil {
//		if os.IsNotExist(err) {
//			log.Fatal("File does not exist.")
//			return err
//		}
//	}
//	// Wait for pending write operations
//	h.WaitGroup.Wait()
//
//	// Lock RW Mutex
//	h.Mutex.Lock()
//
//	// Notify wg start of operation
//	h.WaitGroup.Add(1)
//
//	// Write file
//	err = ioutil.WriteFile("/Users/ijinkus/go/src/paxos/src/api/database/db.json", data, 0666)
//	if err != nil {
//		// If error unlock mutex and notify wg
//		log.Fatal(err)
//		h.Mutex.Unlock()
//		h.WaitGroup.Done()
//		return err
//	}
//
//	// Unlock Mutex
//	h.Mutex.Unlock()
//
//	// Notify wg end of operation
//	h.WaitGroup.Done()
//
//	return nil
//
//}

func (h DbHandler) GetData() models.DB {
	// Wait for pending write operations
	h.WaitGroup.Wait()
	return Data
}

func (h DbHandler) SetData(newTransaction models.Transaction) {
	// Wait for pending write operations
	h.WaitGroup.Wait()

	// Lock RW Mutex
	h.Mutex.Lock()

	// Notify wg start of operation
	h.WaitGroup.Add(1)

	if newTransaction.Type=="credit" {
		Data.Balance.Balance += newTransaction.Ammount
	}
	if newTransaction.Type=="debit" {
		Data.Balance.Balance -= newTransaction.Ammount
	}

	Data.History = append(Data.History, newTransaction)

	// Unlock Mutex
	h.Mutex.Unlock()

	// Notify wg end of operation
	h.WaitGroup.Done()

}