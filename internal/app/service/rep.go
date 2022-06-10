package service

import "example.com/greetings/internal/app/models"

type ServiceRep interface {
	GetStatus() (models.Service, error)
	Clear() error
}
