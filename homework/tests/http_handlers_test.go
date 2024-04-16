//go:build integration
// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"homework/internal/pkg/repository/postgresql"
	"homework/internal/service/service_with_http"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {

	t.Run("smoke test", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		ts := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer ts.Close()

		requestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		// act
		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "{\"ID\":1,\"name\":\"Name\",\"address\":\"Address\",\"phone_number\":\"PhoneNumber\"}", string(respBody))
	})

	t.Run("empty name", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		ts := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer ts.Close()

		requestBody := service_with_http.AddPickupPointRequest{
			Name:        "",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		// act
		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

	t.Run("empty address", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		ts := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer ts.Close()

		requestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "",
			PhoneNumber: "PhoneNumber",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		// act
		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

	t.Run("empty phone number", func(t *testing.T) {
		// arrange

		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		ts := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer ts.Close()

		requestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		// act
		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})
}

func TestUpdate(t *testing.T) {

	t.Run("smoke test", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		tsUpdate := httptest.NewServer(http.HandlerFunc(srv.Update))
		defer tsUpdate.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)

		updateRequestBody := service_with_http.UpdatePickupPointRequest{
			ID:          1,
			Name:        "Updated_Name",
			Address:     "Updated_Address",
			PhoneNumber: "Updated_PhoneNumber",
		}
		updateBodyBytes, _ := json.Marshal(updateRequestBody)
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))

		// act

		updateReq, err := http.NewRequest(http.MethodPut, tsUpdate.URL, bytes.NewBuffer(updateBodyBytes))
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(updateReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "{\"ID\":1,\"name\":\"Updated_Name\",\"address\":\"Updated_Address\",\"phone_number\":\"Updated_PhoneNumber\"}", string(respBody))
	})

	t.Run("empty id", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		tsUpdate := httptest.NewServer(http.HandlerFunc(srv.Update))
		defer tsUpdate.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)

		updateRequestBody := service_with_http.UpdatePickupPointRequest{
			ID:          0,
			Name:        "Updated_Name",
			Address:     "Updated_Address",
			PhoneNumber: "Updated_PhoneNumber",
		}
		updateBodyBytes, _ := json.Marshal(updateRequestBody)
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))

		// act

		updateReq, err := http.NewRequest(http.MethodPut, tsUpdate.URL, bytes.NewBuffer(updateBodyBytes))
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(updateReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

	t.Run("empty name", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		tsUpdate := httptest.NewServer(http.HandlerFunc(srv.Update))
		defer tsUpdate.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)

		updateRequestBody := service_with_http.UpdatePickupPointRequest{
			ID:          1,
			Name:        "",
			Address:     "Updated_Address",
			PhoneNumber: "Updated_PhoneNumber",
		}
		updateBodyBytes, _ := json.Marshal(updateRequestBody)

		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))

		// act

		updateReq, err := http.NewRequest(http.MethodPut, tsUpdate.URL, bytes.NewBuffer(updateBodyBytes))
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(updateReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

	t.Run("empty address", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		tsUpdate := httptest.NewServer(http.HandlerFunc(srv.Update))
		defer tsUpdate.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)

		updateRequestBody := service_with_http.UpdatePickupPointRequest{
			ID:          1,
			Name:        "Updated_Name",
			Address:     "",
			PhoneNumber: "Updated_PhoneNumber",
		}
		updateBodyBytes, _ := json.Marshal(updateRequestBody)
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))

		// act
		updateReq, err := http.NewRequest(http.MethodPut, tsUpdate.URL, bytes.NewBuffer(updateBodyBytes))
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(updateReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

	t.Run("empty phone number", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		tsUpdate := httptest.NewServer(http.HandlerFunc(srv.Update))
		defer tsUpdate.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)

		updateRequestBody := service_with_http.UpdatePickupPointRequest{
			ID:          1,
			Name:        "Updated_Name",
			Address:     "Updated_Address",
			PhoneNumber: "",
		}
		updateBodyBytes, _ := json.Marshal(updateRequestBody)

		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))

		// act
		updateReq, err := http.NewRequest(http.MethodPut, tsUpdate.URL, bytes.NewBuffer(updateBodyBytes))
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(updateReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

	t.Run("not found", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsUpdate := httptest.NewServer(http.HandlerFunc(srv.Update))
		defer tsUpdate.Close()

		updateRequestBody := service_with_http.UpdatePickupPointRequest{
			ID:          1,
			Name:        "Updated_Name",
			Address:     "Updated_Address",
			PhoneNumber: "Updated_PhoneNumber",
		}
		updateBodyBytes, _ := json.Marshal(updateRequestBody)

		// act
		updateReq, err := http.NewRequest(http.MethodPut, tsUpdate.URL, bytes.NewBuffer(updateBodyBytes))
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(updateReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

}

func TestGetByID(t *testing.T) {

	t.Run("smoke test", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		router := mux.NewRouter()
		router.HandleFunc("/{key}", srv.GetByID).Methods(http.MethodGet)
		tsGet := httptest.NewServer(router)
		defer tsGet.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)

		// act
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))
		url := tsGet.URL + "/1"
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "{\"ID\":1,\"name\":\"Name\",\"address\":\"Address\",\"phone_number\":\"PhoneNumber\"}", string(respBody))
	})

	t.Run("not found", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		router := mux.NewRouter()
		router.HandleFunc("/{key}", srv.GetByID).Methods(http.MethodGet)
		tsGet := httptest.NewServer(router)
		defer tsGet.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)

		// act
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))
		url := tsGet.URL + "/12345"
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

}

func TestDelete(t *testing.T) {

	t.Run("smoke test", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		router := mux.NewRouter()
		router.HandleFunc("/{key}", srv.Delete).Methods(http.MethodDelete)
		tsDelete := httptest.NewServer(router)
		defer tsDelete.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)

		// act
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))
		url := tsDelete.URL + "/1"
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

	t.Run("not found", func(t *testing.T) {
		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		router := mux.NewRouter()
		router.HandleFunc("/{key}", srv.Delete).Methods(http.MethodDelete)
		tsDelete := httptest.NewServer(router)
		defer tsDelete.Close()

		createRequestBody := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody)
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))

		// act
		url := tsDelete.URL + "/12345"
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)
		client := &http.Client{}

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		//assert
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "", string(respBody))
	})

}

func TestList(t *testing.T) {

	t.Run("smoke test", func(t *testing.T) {

		// arrange
		db.SetUp(t, "pickup_points")
		defer db.TearDown()

		repo := postgresql.NewPickupPoints(db.DB)
		srv := service_with_http.Server1{repo}
		tsCreate := httptest.NewServer(http.HandlerFunc(srv.Create))
		defer tsCreate.Close()

		tsList := httptest.NewServer(http.HandlerFunc(srv.List))
		defer tsList.Close()

		createRequestBody1 := service_with_http.AddPickupPointRequest{
			Name:        "Name",
			Address:     "Address",
			PhoneNumber: "PhoneNumber",
		}
		createBodyBytes, _ := json.Marshal(createRequestBody1)
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))

		createRequestBody2 := service_with_http.AddPickupPointRequest{
			Name:        "Name2",
			Address:     "Address2",
			PhoneNumber: "PhoneNumber2",
		}
		createBodyBytes, _ = json.Marshal(createRequestBody2)
		http.Post(tsCreate.URL, "application/json", bytes.NewBuffer(createBodyBytes))

		// act
		req, err := http.NewRequest(http.MethodGet, tsList.URL, nil)
		require.NoError(t, err)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		//assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		assert.Equal(t, "[{\"ID\":1,\"name\":\"Name\",\"address\":\"Address\",\"phone_number\":\"PhoneNumber\"},{\"ID\":2,\"name\":\"Name2\",\"address\":\"Address2\",\"phone_number\":\"PhoneNumber2\"}]", string(respBody))
	})

}
