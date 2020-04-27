package database

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type DbInterface interface {
	ReadFile() ([]byte, error)
	WriteFile([]byte) error
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

func (h DbHandler) ReadFile() ([]byte, error) {
	// Check if file exists
	_, err := os.Stat("/Users/ijinkus/go/src/paxos/src/api/database/db.txt")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
			return nil, err
		}
	}
	// Wait for pending write operations
	h.WaitGroup.Wait()

	data, err := ioutil.ReadFile("/Users/ijinkus/go/src/paxos/src/api/database/db.txt")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return data, nil
}

func (h DbHandler) WriteFile(data []byte) error {
	// Check if file exists
	_, err := os.Stat("/Users/ijinkus/go/src/paxos/src/api/database/db.txt")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
			return err
		}
	}
	// Wait for pending write operations
	h.WaitGroup.Wait()

	// Lock RW Mutex
	h.Mutex.Lock()

	// Notify wg start of operation
	h.WaitGroup.Add(1)

	// Write file
	err = ioutil.WriteFile("/Users/ijinkus/go/src/paxos/src/api/database/db.txt", data, 0666)
	if err != nil {
		// If error unlock mutex and notify wg
		log.Fatal(err)
		h.Mutex.Unlock()
		h.WaitGroup.Done()
		return err
	}

	// Unlock Mutex
	h.Mutex.Unlock()

	// Notify wg end of operation
	h.WaitGroup.Done()

	return nil

}