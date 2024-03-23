package main

import (
	"Homework-1/internal/pkg/db"
	"Homework-1/internal/pkg/repository"
	"Homework-1/internal/pkg/repository/postgresql"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const port = ":9000"
const queryParamKey = "key"

type server1 struct {
	repo *postgresql.PickupPointRepo
}

type addPickupPointRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type addPickupPointResponse struct {
	ID          int64  `json:"ID"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type updatePickupPointRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	pickupPointsRepo := postgresql.NewPickupPoints(database)
	implemetation := server1{repo: pickupPointsRepo}

	http.Handle("/", createRouter(implemetation))
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func createRouter(implemetation server1) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/pickup_point", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implemetation.Create(w, req)
		case http.MethodPut:
			implemetation.Update(w, req)
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc(fmt.Sprintf("/pickup_point/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			fmt.Println("Getbyid")
			implemetation.GetByID(w, req)
		case http.MethodDelete:
			implemetation.Delete(w, req)
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc("/pickup_point/list", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			fmt.Println("Getbyid")
			implemetation.List(w, req)
		default:
			fmt.Println("error")
		}
	})

	return router
}

func (s *server1) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm addPickupPointRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pickupPointRepo := &repository.PickupPoint{
		Name:        unm.Name,
		Address:     unm.Address,
		PhoneNumber: unm.PhoneNumber,
	}
	id, err := s.repo.Add(req.Context(), pickupPointRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &addPickupPointResponse{
		ID:          id,
		Name:        pickupPointRepo.Name,
		Address:     pickupPointRepo.Address,
		PhoneNumber: pickupPointRepo.PhoneNumber,
	}
	pickupPointJson, _ := json.Marshal(resp)
	w.Write(pickupPointJson)
}

func (s *server1) Update(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm updatePickupPointRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := unm.ID
	pickupPointRepo := &repository.PickupPoint{
		Name:        unm.Name,
		Address:     unm.Address,
		PhoneNumber: unm.PhoneNumber,
	}
	id, err = s.repo.Update(req.Context(), id, pickupPointRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &addPickupPointResponse{
		ID:          id,
		Name:        pickupPointRepo.Name,
		Address:     pickupPointRepo.Address,
		PhoneNumber: pickupPointRepo.PhoneNumber,
	}
	pickupPointJson, _ := json.Marshal(resp)
	w.Write(pickupPointJson)
}

func (s *server1) GetByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		fmt.Println(ok)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	point, err := s.repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pointJson, _ := json.Marshal(point)
	w.Write(pointJson)
}

func (s *server1) Delete(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		fmt.Println(ok)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.Delete(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *server1) List(w http.ResponseWriter, req *http.Request) {
	points, err := s.repo.List(req.Context())
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pointsJson, _ := json.Marshal(points)
	w.Write(pointsJson)
}
