package card

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrTypeDoesNotExist   = errors.New("card type does not exist")
	ErrIssuerDoesNotExist = errors.New("card issuer does not exist")
	ErrUserDoesNotExist   = errors.New("user does not exist")
	ErrNoBaseCard         = errors.New("user dont have base card")
)

type Card struct {
	Issuer, Type       string
	Id, UserId, Number int64
}

type Service struct {
	mu    sync.RWMutex
	Cards []*Card
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) All(context.Context) []*Card {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Cards
}

func setNumber(num int64) int64 {
	num += 1
	return num
}

func (s *Service) Add(userId int64, typeCard, issuerCard string) (*Card, error) {
	err := getIssuerCard(issuerCard)
	if err != nil {
		return &Card{}, err
	}

	err = getTypeCard(typeCard)
	if err != nil {
		return &Card{}, err
	}

	number, err := s.getBaseCard(userId)
	if err != nil && typeCard != "basic" {
		return &Card{}, err
	}

	card := &Card{
		Id:     setNumber(number),
		UserId: userId,
		Number: setNumber(number),
		Type:   typeCard,
		Issuer: issuerCard,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Cards = append(s.Cards, card)
	return card, nil
}

func getIssuerCard(issuerCard string) error {
	issuers := []string{"Visa", "Maestro", "MasterCard"}
	for _, v := range issuers {
		if v == issuerCard {
			return nil
		}
	}
	return ErrIssuerDoesNotExist
}

func getTypeCard(typeCard string) error {
	types := []string{"basic", "additional", "virtual"}
	for _, value := range types {
		if value == typeCard {
			return nil
		}
	}
	return ErrTypeDoesNotExist
}

func (s *Service) getBaseCard(userId int64) (number int64, err error) {
	for _, value := range s.Cards {
		if value.UserId == userId {
			return value.Number, nil
		}
	}
	return 0, ErrNoBaseCard
}
