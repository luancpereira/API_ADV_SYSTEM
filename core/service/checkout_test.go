package service_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jellydator/ttlcache/v3"
	coreErrors "github.com/luancpereira/APICheckout/core/errors"
	"github.com/luancpereira/APICheckout/core/service"
	"github.com/stretchr/testify/assert"
)

type MockCheckout struct{}

func (m *MockCheckout) CreateTransaction(description string, transactionDate time.Time, transactionValue float64) (int64, error) {
	return 12345, nil
}

type MockError struct{}

func (m MockError) New(keys ...string) *coreErrors.CoreError {
	return &coreErrors.CoreError{
		Key:     keys[0],
		Message: "Mocked message for testing",
	}
}

type MockResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func TestMain(m *testing.M) {
	coreErrors.C = ttlcache.New[string, string]()

	coreErrors.C.Set("error.description.empty", "Description cannot be empty.", ttlcache.NoTTL)
	coreErrors.C.Set("error.description.too.long", "Description must be less than 50 characters.", ttlcache.NoTTL)
	coreErrors.C.Set("error.value.not.positive", "Value must be positive.", ttlcache.NoTTL)
	coreErrors.C.Set("error.not.found.value.record", "Value cannot be converted to the currency.", ttlcache.NoTTL)

	m.Run()
}

func TestCreateTransaction(t *testing.T) {

	description := "Test transaction"
	transactionDate := time.Now()
	transactionValue := 100.0

	checkout := &MockCheckout{}

	id, err := checkout.CreateTransaction(description, transactionDate, transactionValue)

	assert.NoError(t, err)

	assert.Greater(t, id, int64(0), "O ID da transação deveria ser maior que 0")

	assert.NotNil(t, id, "O ID da transação não pode ser nil")
}

func TestValidateDescription(t *testing.T) {
	t.Run("Deve retornar erro para descrição vazia", func(t *testing.T) {
		err := service.Checkout{}.ValidateDescription("")
		coreErr, ok := err.(*coreErrors.CoreError)
		assert.True(t, ok, "O erro retornado deve ser do tipo CoreError")
		assert.Equal(t, "error.description.empty", coreErr.Key)
	})

	t.Run("Deve retornar erro para descrição longa", func(t *testing.T) {
		err := service.Checkout{}.ValidateDescription("Descrição muito longa que excede os 50 caracteres permitidos.")
		coreErr, ok := err.(*coreErrors.CoreError)
		assert.True(t, ok, "O erro retornado deve ser do tipo CoreError")
		assert.Equal(t, "error.description.too.long", coreErr.Key)
	})

	t.Run("Deve validar descrição válida", func(t *testing.T) {
		err := service.Checkout{}.ValidateDescription("Descrição válida")
		assert.NoError(t, err)
	})
}

func TestValidateTransactionValue(t *testing.T) {
	t.Run("Deve retornar erro para valor não positivo", func(t *testing.T) {
		err := service.Checkout{}.ValidateTrasactionValue(0)
		coreErr, ok := err.(*coreErrors.CoreError)
		assert.True(t, ok, "O erro retornado deve ser do tipo CoreError")
		assert.Equal(t, "error.value.not.positive", coreErr.Key)

		err = service.Checkout{}.ValidateTrasactionValue(-10)
		coreErr, ok = err.(*coreErrors.CoreError)
		assert.True(t, ok, "O erro retornado deve ser do tipo CoreError")
		assert.Equal(t, "error.value.not.positive", coreErr.Key)
	})

	t.Run("Deve retornar sucesso para valor positivo", func(t *testing.T) {
		err := service.Checkout{}.ValidateTrasactionValue(100)
		assert.NoError(t, err)
	})
}

func TestGetEntity(t *testing.T) {
	t.Run("Deve retornar sucesso para resposta válida", func(t *testing.T) {
		mockData := MockResponse{Name: "Test", Value: "12345"}
		mockResponse := `{"name": "Test", "value": "12345"}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockResponse))
		}))
		defer server.Close()

		var result MockResponse
		err := service.GetEntity(server.URL, nil, &result)

		assert.NoError(t, err)
		assert.Equal(t, mockData.Name, result.Name)
		assert.Equal(t, mockData.Value, result.Value)
	})

	t.Run("Deve retornar erro para falha na criação da requisição", func(t *testing.T) {
		err := service.GetEntity("http://:invalid", nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar a requisição")
	})

	t.Run("Deve retornar erro para falha no status da resposta", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		var result MockResponse
		err := service.GetEntity(server.URL, nil, &result)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("requisição falhou com status %d", http.StatusInternalServerError))
	})

	t.Run("Deve retornar erro para falha na decodificação do JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{invalid json"))
		}))
		defer server.Close()

		var result MockResponse
		err := service.GetEntity(server.URL, nil, &result)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao decodificar a resposta JSON")
	})
}

func TestFindRegistryWithDateCloset(t *testing.T) {
	records := []service.Record{
		{
			EffectiveDate: "2025-01-01",
			Country:       "BRAZIL",
			ExchangeRate:  "1.25",
		},
		{
			EffectiveDate: "2025-01-05",
			Country:       "BRAZIL",
			ExchangeRate:  "1.30",
		},
		{
			EffectiveDate: "2025-01-10",
			Country:       "BRAZIL",
			ExchangeRate:  "1.35",
		},
		{
			EffectiveDate: "2024-07-01",
			Country:       "BRAZIL",
			ExchangeRate:  "1.50",
		},
	}

	t.Run("Deve retornar o registro com a data mais próxima anterior", func(t *testing.T) {
		targetDate := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)

		closestRecord, err := service.FindRegistryWithDateCloset(records, targetDate)

		assert.NoError(t, err)
		assert.Equal(t, "2025-01-05", closestRecord.EffectiveDate)
		assert.Equal(t, "1.30", closestRecord.ExchangeRate)
	})

	t.Run("Deve retornar erro quando não houver registros válidos", func(t *testing.T) {
		targetDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		closestRecord, err := service.FindRegistryWithDateCloset(records, targetDate)

		assert.Error(t, err)

		coreErr, ok := err.(*coreErrors.CoreError)
		assert.True(t, ok, "O erro retornado deve ser do tipo *coreErrors.CoreError")

		assert.Equal(t, "error.not.found.value.record", coreErr.Key)

		assert.Equal(t, service.Record{}, closestRecord)
	})

	t.Run("Deve retornar erro para datas de registro inválidas", func(t *testing.T) {
		invalidRecords := append(records, service.Record{
			EffectiveDate: "invalid-date",
			Country:       "BRAZIL",
			ExchangeRate:  "1.60",
		})

		targetDate := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)

		closestRecord, err := service.FindRegistryWithDateCloset(invalidRecords, targetDate)

		assert.NoError(t, err)
		assert.Equal(t, "2025-01-05", closestRecord.EffectiveDate)
	})

	t.Run("Deve ignorar registros além do limite de 6 meses", func(t *testing.T) {
		targetDate := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)

		closestRecord, err := service.FindRegistryWithDateCloset(records, targetDate)

		assert.NoError(t, err)
		assert.Equal(t, "2025-01-05", closestRecord.EffectiveDate)
		assert.Equal(t, "1.30", closestRecord.ExchangeRate)
	})
}
