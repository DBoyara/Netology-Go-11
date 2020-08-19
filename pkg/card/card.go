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

type UserCards []*Card

type UserID int64

type Card struct {
	Id     int64
	Issuer string
	Type   string
	Number int64
}

type Service struct {
	mu     sync.RWMutex
	Cards  map[UserID]UserCards
	lastID int64
}

func NewService() *Service {
	return &Service{
		mu:     sync.RWMutex{},
		Cards:  map[UserID]UserCards{},
		lastID: 0,
	}
}

func (s *Service) All(id UserID) (UserCards, error) {
	cards, err := s.getCardsByUserID(id)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (c UserCards) nextCardNumber() int64 {
	if len(c) == 0 {
		return 1
	}
	i := c[len(c)-1]
	return setNumber(i.Number)
}

func setNumber(num int64) int64 {
	num += 1
	return num
}

func (s *Service) getCardsByUserID(id UserID) (UserCards, error) {
	v, ok := s.Cards[id]
	if !ok {
		return nil, ErrNoBaseCard
	}
	return v, nil
}

func (s *Service) Add(userId UserID, typeCard, issuerCard string) (*Card, error) {

	cards, err := s.getCardsByUserID(userId)
	if err != nil && typeCard != "base" {
		return nil, err
	}

	err = getIssuerCard(issuerCard)
	if err != nil {
		return nil, err
	}

	err = getTypeCard(typeCard)
	if err != nil {
		return nil, err
	}

	s.lastID = cards.nextCardNumber()

	newCard := &Card{
		Id:     s.lastID,
		Issuer: issuerCard,
		Type:   typeCard,
		Number: s.lastID,
	}
	s.Cards[userId] = append(s.Cards[userId], newCard)

	return newCard, nil
}

func getIssuerCard(issuerCard string) error {
	issuers := map[string]struct{}{
		"Visa":       {},
		"Maestro":    {},
		"MasterCard": {},
	}

	if _, ok := issuers[issuerCard]; !ok {
		return ErrIssuerDoesNotExist
	}

	return nil
}

func getTypeCard(typeCard string) error {
	types := map[string]struct{}{
		"base":       {},
		"additional": {},
		"virtual":    {},
	}

	if _, ok := types[typeCard]; !ok {
		return ErrTypeDoesNotExist
	}

	return nil

}
