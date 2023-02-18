package service

import "github.com/Hymiside/evraz-api/pkg/repository"

type Evraz interface {
	Register(nameConn string) WrapChan
	Unregister(nameConn string)
	Consume(d map[string]interface{})
}

type Service struct {
	Evraz
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Evraz: NewEvrazService(repo),
	}
}
