package service_with_http

import (
	"Homework-1/internal/infrastructure/answer"
	"Homework-1/internal/infrastructure/info"
	"Homework-1/internal/infrastructure/kafka"
	"Homework-1/internal/pkg/db"
	"Homework-1/internal/pkg/repository"
	"Homework-1/internal/pkg/repository/postgresql"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	securePort   = ":9000"
	insecurePort = ":9001"
)

const queryParamKey = "key"

type Server1 struct {
	Repo     repository.PickupPointRepo
	sender   answer.Sender
	receiver info.Receiver
}

type infoMessage struct {
	Timestamp time.Time
	Method    string
	Raw       string
}

type AddPickupPointRequest struct {
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

type UpdatePickupPointRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

func Secure() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	pickupPointsRepo := postgresql.NewPickupPoints(database)
	implemetation := Server1{Repo: pickupPointsRepo}
	mx := http.NewServeMux()

	broker, _ := os.LookupEnv("BROKER")
	kafkaProducer, err := kafka.NewProducer([]string{broker})
	sender := answer.NewKafkaSender(kafkaProducer, "info")

	kafkaConsumer, err := kafka.NewConsumer([]string{broker})
	handlers := map[string]info.HandleFunc{
		"payments": func(message *sarama.ConsumerMessage) {
			pm := infoMessage{}
			err = json.Unmarshal(message.Value, &pm)
			if err != nil {
				fmt.Println("Consumer error", err)
			}

			fmt.Println("Received Key: ", string(message.Key), " Value: ", pm)
		},
	}
	infos := info.NewService(info.NewReceiver(kafkaConsumer, handlers))
	infos.StartConsume("info")

	mx.Handle("/", authMiddleware(createRouter(implemetation)))
	if err := http.ListenAndServeTLS(securePort, "./server.crt", "./server.key", mx); err != nil {
		log.Fatal(err)
	}

}

func Insecure() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	pickupPointsRepo := postgresql.NewPickupPoints(database)
	implemetation := Server1{Repo: pickupPointsRepo}
	mx := http.NewServeMux()
	mx.Handle("/", authMiddleware(createRouter(implemetation)))
	if err := http.ListenAndServe(insecurePort, mx); err != nil {
		log.Fatal(err)
	}

}

func createRouter(implemetation Server1) *mux.Router {
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

func (s *Server1) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var unm AddPickupPointRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.validateCreate(req.Context(), unm)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pickupPointJson, status, err := s.create(req.Context(), unm)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(status)
	w.Write(pickupPointJson)
}

func (s *Server1) validateCreate(ctx context.Context, unm AddPickupPointRequest) error {
	if unm.Name == "" {
		return errors.New("Name field is empty")
	}
	if unm.Address == "" {
		return errors.New("Address field is empty")
	}
	if unm.PhoneNumber == "" {
		return errors.New("PhoneNumber field is empty")
	}
	return nil
}

func (s *Server1) create(ctx context.Context, unm AddPickupPointRequest) ([]byte, int, error) {
	pickupPoint := &repository.PickupPoint{
		Name:        unm.Name,
		Address:     unm.Address,
		PhoneNumber: unm.PhoneNumber,
	}
	id, err := s.Repo.Add(ctx, pickupPoint)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	resp := addPickupPointResponse{
		ID:          id,
		Name:        pickupPoint.Name,
		Address:     pickupPoint.Address,
		PhoneNumber: pickupPoint.PhoneNumber,
	}
	pickupPointJson, _ := json.Marshal(resp)

	return pickupPointJson, http.StatusOK, nil
}

func (s *Server1) Update(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var unm UpdatePickupPointRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.validateUpdate(req.Context(), unm)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pickupPointJson, status, err := s.update(req.Context(), unm)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(status)
	w.Write(pickupPointJson)
}

func (s *Server1) validateUpdate(ctx context.Context, unm UpdatePickupPointRequest) error {
	if unm.ID == 0 {
		return errors.New("ID field is empty")
	}
	if unm.Name == "" {
		return errors.New("Name field is empty")
	}
	if unm.Address == "" {
		return errors.New("Address field is empty")
	}
	if unm.PhoneNumber == "" {
		return errors.New("PhoneNumber field is empty")
	}
	return nil
}

func (s *Server1) update(ctx context.Context, unm UpdatePickupPointRequest) ([]byte, int, error) {
	id := unm.ID
	pickupPointRepo := &repository.PickupPoint{
		Name:        unm.Name,
		Address:     unm.Address,
		PhoneNumber: unm.PhoneNumber,
	}
	err := s.Repo.Update(ctx, id, pickupPointRepo)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("Could not find object with id =%d", id))
		}
		return nil, http.StatusInternalServerError, err
	}

	resp := &addPickupPointResponse{
		ID:          id,
		Name:        pickupPointRepo.Name,
		Address:     pickupPointRepo.Address,
		PhoneNumber: pickupPointRepo.PhoneNumber,
	}
	pickupPointJson, _ := json.Marshal(resp)

	return pickupPointJson, http.StatusOK, nil
}

func (s *Server1) GetByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		log.Println("could not parse object id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pointJson, status, err := s.get(req.Context(), keyInt)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(status)
	w.Write(pointJson)
}

func (s *Server1) get(ctx context.Context, key int64) ([]byte, int, error) {
	point, err := s.Repo.GetByID(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound, errors.New(fmt.Sprintf("Could not find object with id =%d", key))
		}
		return nil, http.StatusInternalServerError, err
	}
	pointJson, _ := json.Marshal(point)
	return pointJson, http.StatusOK, nil
}

func (s *Server1) Delete(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		log.Println("Could not parse object id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status, err := s.delete(req.Context(), keyInt)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(status)
}

func (s *Server1) delete(ctx context.Context, keyInt int64) (int, error) {
	err := s.Repo.Delete(ctx, keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			log.Println()
			return http.StatusNotFound, errors.New(fmt.Sprintf("Could not find object with id =%d", keyInt))
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *Server1) List(w http.ResponseWriter, req *http.Request) {
	pointsJson, status, err := s.list(req.Context())
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(status)
	w.Write(pointsJson)
}

func (s *Server1) list(ctx context.Context) ([]byte, int, error) {
	points, err := s.Repo.List(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	pointsJson, _ := json.Marshal(points)
	return pointsJson, http.StatusOK, nil
}

func authMiddleware(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		true_username, username_exists := os.LookupEnv("USER")
		true_password, password_exists := os.LookupEnv("PASSWORD")
		if username_exists && password_exists {
			username, password, ok := req.BasicAuth()
			if !ok || username != true_username || password != true_password {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		if req.Method == http.MethodPost || req.Method == http.MethodPut {
			body, err := io.ReadAll(req.Body)
			req.Body.Close() //  must close
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			var unm UpdatePickupPointRequest
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
