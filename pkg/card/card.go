package card

import (
	"errors"
	"sync"
)

var (
	ErrTypeDoesNotExist   = errors.New("card type does not exist")
	ErrIssuerDoesNotExist = errors.New("card issuer does not exist")
	ErrUserDoesNotExist   = errors.New("user does not exist")
	ErrNoBaseCard         = errors.New("user dont have base card")
	ErrNotSpecifiedUserId = errors.New("user id unspecified ")
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

func (s *Service) All() []*Card {
	return s.Cards
}

func setNumber(num int64) int64 {
	num += 1
	return num
}

func (s *Service) Add(userId int64, typeCard, idIssuerCard string) (*Card, error) {
	err := getIssuerCard(idIssuerCard)
	if err != nil {
		return &Card{}, err
	}

	err = getTypeCard(typeCard)
	if err != nil {
		return &Card{}, err
	}

	number, err := s.getBaseCard(userId)
	if err != nil && typeCard != "base" {
		return &Card{}, err
	}

	card := &Card{
		Id:     setNumber(number),
		UserId: userId,
		Number: setNumber(number),
		Type:   typeCard,
		Issuer: idIssuerCard,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Cards = append(s.Cards, card)
	return card, nil
}

func getIssuerCard(issuerCard string) error {
	issuers := map[string]struct{}{
		"Visa": {},
		"Maestro": {},
		"MasterCard": {},
	}

	if _, ok := issuers[issuerCard]; !ok {
		return ErrIssuerDoesNotExist
	}

	return nil
}

func getTypeCard(typeCard string) error {
	types := map[string]struct{}{
		"basic": {},
		"additional": {},
		"virtual": {},
	}

	if _, ok := types[typeCard]; !ok {
		return ErrTypeDoesNotExist
	}

	return nil

}

func (s *Service) getBaseCard(userId int64) (number int64, err error) {
	for _, value := range s.Cards {
		if value.UserId == userId {
			return value.Number, nil
		}
	}
	return 0, ErrNoBaseCard
}
