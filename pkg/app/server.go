package app

import (
	"encoding/json"
	"fmt"
	"github.com/DBoyara/Netology-Go-11/pkg/app/dto"
	"github.com/DBoyara/Netology-Go-11/pkg/card"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Server struct {
	cardSvc *card.Service
	mux     *http.ServeMux
}

func NewServer(cardSvc *card.Service, mux *http.ServeMux) *Server {
	return &Server{cardSvc: cardSvc, mux: mux}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/getCards", s.getCards)
	s.mux.HandleFunc("/addCard", s.addCard)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	userId, ok := url.Parse(r.URL.Query().Get("userId"))
	intUserId, err := strconv.Atoi(fmt.Sprintf("%v", userId))
	
	if ok != nil || intUserId == 0 || err != nil {
		dtos := dto.CardErrDTO{Err: card.ErrUserDoesNotExist.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	cards := s.cardSvc.All(r.Context())
	var dtos []*dto.CardDTO

	for _, c := range cards {
		if int64(intUserId) == c.UserId {
			dtos = append(
				dtos,
				&dto.CardDTO{
					Id:     c.Id,
					UserId: c.UserId,
					Number: c.Number,
					Type:   c.Type,
					Issuer: c.Issuer,
				})
		}
	}
	jsonResponse(w, r, dtos)
}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}

	params := &dto.CardDTO{}
	err = json.Unmarshal(body, params)
	if err != nil {
		log.Print(err)
	}
	newCard, err := s.cardSvc.Add(params.UserId, params.Type, params.Issuer)

	if err != nil {
		dtos := dto.CardErrDTO{Err: err.Error()}
		jsonResponse(w, r, dtos)
		return
	}

	dtos := []*dto.CardDTO{}
	dtos = append(dtos,
		&dto.CardDTO{
			Id:     newCard.Id,
			UserId: newCard.UserId,
			Number: newCard.Number,
			Type:   newCard.Type,
			Issuer: newCard.Issuer,
		})

	jsonResponse(w, r, dtos)
}

func jsonResponse(w http.ResponseWriter, r *http.Request, dtos... interface{}) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
		return
	}
}
