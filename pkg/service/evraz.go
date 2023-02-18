package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/Hymiside/evraz-api/pkg/repository"
)

type Data struct {
	t      time.Time
	values map[string]interface{}
}

type WrapChan struct {
	Ch       chan Data
	IsClosed *bool
}

type EvrazService struct {
	repo    repository.Evraz
	mu      *sync.RWMutex
	storage map[string]WrapChan
}

func NewEvrazService(repo repository.Evraz) *EvrazService {
	return &EvrazService{repo: repo, mu: &sync.RWMutex{}, storage: make(map[string]WrapChan)}
}

func (e *EvrazService) Register(nameConn string) WrapChan {
	f := false
	wc := WrapChan{Ch: make(chan Data), IsClosed: &f}

	fmt.Println("lol")

	e.mu.Lock()
	defer e.mu.Unlock()

	e.storage[nameConn] = wc
	return wc
}

func (e *EvrazService) Unregister(nameConn string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	wc, ok := e.storage[nameConn]
	if !ok {
		return
	}
	*wc.IsClosed = true
}

func (e *EvrazService) Consume(d map[string]interface{}) {
	//fmt.Println(d)

	e.mu.Lock()
	defer e.mu.Unlock()

	//t := fmt.Sprintf("%v", d["moment"])
	t := time.Now()
	delete(d, "moment")

	dataExgauster := Data{t: t, values: d}

	for key, v := range e.storage {
		if v.IsClosed == nil || *v.IsClosed {
			delete(e.storage, key)
		}
		func(c chan Data) {
			c <- dataExgauster
		}(v.Ch)

	}
	// Записать в БД

}
