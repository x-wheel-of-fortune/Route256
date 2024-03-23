package main

import (
	"Homework-1/internal/pkg/db"
	"Homework-1/internal/pkg/repository"
	"Homework-1/internal/pkg/repository/postgresql"
	"bytes"
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

	http.Handle("/", authMiddleware(createRouter(implemetation)))
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
			fmt.Println("This route does not support", req.Method, "method.")
		}
	})

	router.HandleFunc(fmt.Sprintf("/pickup_point/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.GetByID(w, req)
		case http.MethodDelete:
			implemetation.Delete(w, req)
		default:
			fmt.Println("This route does not support", req.Method, "method.")
		}
	})

	router.HandleFunc("/pickup_point/list", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.List(w, req)
		default:
			fmt.Println("This route does not support", req.Method, "method.")
		}
	})

	return router
}

func (s *server1) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var unm addPickupPointRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if unm.Name == "" {
		log.Println("Name field is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if unm.Address == "" {
		log.Println("Address field is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if unm.PhoneNumber == "" {
		log.Println("PhoneNumber field is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pickupPointRepo := &repository.PickupPoint{
		Name:        unm.Name,
		Address:     unm.Address,
		PhoneNumber: unm.PhoneNumber,
	}
	id, err := s.repo.Add(req.Context(), pickupPointRepo)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var unm updatePickupPointRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if unm.ID == 0 {
		log.Println("ID field is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if unm.Name == "" {
		log.Println("Name field is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if unm.Address == "" {
		log.Println("Address field is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if unm.PhoneNumber == "" {
		log.Println("PhoneNumber field is empty")
		w.WriteHeader(http.StatusBadRequest)
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
		log.Println(err)
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
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	point, err := s.repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			log.Println("Could not find object with id =", keyInt)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
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
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.repo.Delete(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			log.Println("Could not find object with id =", keyInt)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *server1) List(w http.ResponseWriter, req *http.Request) {
	points, err := s.repo.List(req.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pointsJson, _ := json.Marshal(points)
	w.Write(pointsJson)
}

func authMiddleware(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//username, password, ok := req.BasicAuth()
		//if ok {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	return
		//}
		//fmt.Println(username, password)

		if req.Method == http.MethodPost || req.Method == http.MethodPut {
			body, err := io.ReadAll(req.Body)
			req.Body.Close() //  must close
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			var unm updatePickupPointRequest
			if err = json.Unmarshal(body, &unm); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Printf("Method: %s, body: %+v\n", req.Method, unm)
		} else if req.Method == http.MethodDelete {
			fmt.Printf("Method: %s, to_be_deleted: %s\n", req.Method, req.URL)
		}

		handler.ServeHTTP(w, req)

	}
}
