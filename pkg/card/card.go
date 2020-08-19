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
	return &Service{}
}

func (s *Service) All(id UserID) (UserCards, error) {
	cards, err := s.getCardsByUserID(id)
	if err != nil {
		return UserCards{}, err
	}
	return cards, nil
}

// раз мы работаем с упорядоченным массивом, то можно просто взять последний элемент с его номером
func (c UserCards) nextCardNumber() int64 {
	if len(c) == 0 {
		return 0
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
		return UserCards{}, ErrNoBaseCard
	}
	return v, nil
}

func (s *Service) Add(userId int64, typeCard, issuerCard string) (*Card, error) {

	cards, err := s.getCardsByUserID(UserID(userId))
	if err != nil && typeCard != "base" {
		return &Card{}, err
	}

	err = getIssuerCard(issuerCard)
	if err != nil {
		return &Card{}, err
	}

	err = getTypeCard(typeCard)
	if err != nil {
		return &Card{}, err
	}

	s.lastID = cards.nextCardNumber()

	newCard := &Card{
		Id:     s.lastID,
		Issuer: issuerCard,
		Type:   typeCard,
		Number: s.lastID,
	}
	cards = append(cards, newCard)

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
