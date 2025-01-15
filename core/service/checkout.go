package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/luancpereira/APICheckout/core/database"
	"github.com/luancpereira/APICheckout/core/database/sqlc"
	coreError "github.com/luancpereira/APICheckout/core/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Checkout struct{}

/*****
funcs for creations
******/

func (c Checkout) CreateTransaction(description string, transaction_date time.Time, transaction_value float64) (ID int64, err error) {

	err = c.ValidateDescription(description)
	if err != nil {
		return
	}

	err = c.ValidateTrasactionValue(transaction_value)
	if err != nil {
		return
	}

	if math.Round(transaction_value*100) != transaction_value*100 {
		transaction_value = math.Round(transaction_value*100) / 100
	}

	params := sqlc.InsertTransactionParams{
		Description:      description,
		TransactionDate:  transaction_date,
		TransactionValue: transaction_value,
	}

	ID, err = database.DB_QUERIER.InsertTransaction(context.Background(), params)
	if err != nil {
		err = database.Utils{}.CoreErrorDatabase(err)
		return
	}

	return
}

/*****
funcs for creations
******/

/*****
funcs for gets
******/

func (Checkout) GetByID(transactionID int64, country string) (transaction TransactionDetail, err error) {
	transactionDetail, err := database.DB_QUERIER.SelectTransactionByID(context.Background(), transactionID)
	if err != nil {
		err = database.Utils{}.CoreErrorDatabase(err)
		return
	}

	exchangeRate, err := getExchangeRate(transactionDetail.TransactionDate, country)
	if err != nil {
		return
	}

	transaction = TransactionDetail{
		SelectTransactionByIDRow:                transactionDetail,
		ExchangeRate:                            math.Round(exchangeRate*100) / 100,
		TransactionValueConvertedToWishCurrency: math.Round(transactionDetail.TransactionValue*exchangeRate*100) / 100,
	}

	return
}

func (Checkout) GetList(filters map[string]string, limit, offset int64, country string) (models []TransactionDetailList, total int64, err error) {
	params := sqlc.SelectTransactionsParams{
		Column1:         limit,
		Column2:         offset,
		TransactionDate: filters["transaction_date"],
	}

	transactions, err := database.DB_QUERIER.SelectTransactions(context.Background(), params)
	if err != nil {
		err = database.Utils{}.CoreErrorDatabase(err)
		return
	}

	parsedDate, _ := time.Parse("2006-01-02", filters["transaction_date"])

	exchangeRate, err := getExchangeRate(parsedDate, country)

	if err != nil {
		return
	}

	var transactionDetailList []TransactionDetailList

	for _, transaction := range transactions {

		transactionDetail := TransactionDetailList{
			SelectTransactionsRow:                   transaction,
			ExchangeRate:                            math.Round(exchangeRate*100) / 100,
			TransactionValueConvertedToWishCurrency: math.Round(transaction.TransactionValue*exchangeRate*100) / 100,
		}

		transactionDetailList = append(transactionDetailList, transactionDetail)
	}

	models = transactionDetailList

	total, err = database.DB_QUERIER.SelectTransactionsTotal(context.Background(), filters["transaction_date"])
	if err != nil {
		err = database.Utils{}.CoreErrorDatabase(err)
		return
	}

	return
}

func getExchangeRate(transactionDate time.Time, country string) (float64, error) {
	formattedDate := transactionDate.Format("2006-01-02")
	url := "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?filter=country:eq:" + CapitalizeFirstLetter(country) + ",effective_date:lte:" + formattedDate

	var response Response
	err := GetEntity(url, map[string]string{}, &response)
	if err != nil {
		return 0, err
	}

	closestRecord, err := FindRegistryWithDateCloset(response.Data, transactionDate)
	if err != nil {
		return 0, err
	}

	exchangeRate, err := strconv.ParseFloat(closestRecord.ExchangeRate, 64)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter ExchangeRate para float64: %w", err)
	}

	return exchangeRate, nil
}

/*****
funcs for gets
******/

/*****
funcs for validations
******/

func (Checkout) ValidateDescription(description string) (err error) {
	if len(description) == 0 {
		err = coreError.New("error.description.empty")
		return
	}

	if len(description) > 50 {
		err = coreError.New("error.description.too.long")
		return
	}

	return
}

func (Checkout) ValidateTrasactionValue(value float64) (err error) {
	if value <= 0 {
		err = coreError.New("error.value.not.positive")
		return
	}

	return
}

func FindRegistryWithDateCloset(records []Record, targetDate time.Time) (closestRecord Record, err error) {
	var minDiff time.Duration = time.Duration(math.MaxInt64)

	const maxDuration = 182 * 24 * time.Hour

	for _, record := range records {
		recordDate, err := time.Parse("2006-01-02", record.EffectiveDate)
		if err != nil {
			continue
		}

		diff := recordDate.Sub(targetDate)
		if diff < 0 {
			diff = -diff
		}

		if diff < minDiff && recordDate.Before(targetDate) && diff <= maxDuration {
			minDiff = diff
			closestRecord = record
		}
	}

	if minDiff == time.Duration(math.MaxInt64) {
		err = coreError.New("error.not.found.value.record")
	}

	return
}

/*****
funcs for validations
******/

/*****
other funcs
******/

type TransactionDetail struct {
	sqlc.SelectTransactionByIDRow
	ExchangeRate                            float64
	TransactionValueConvertedToWishCurrency float64
}

type TransactionDetailList struct {
	sqlc.SelectTransactionsRow
	ExchangeRate                            float64
	TransactionValueConvertedToWishCurrency float64
}

type Record struct {
	RecordDate            string `json:"record_date"`
	Country               string `json:"country"`
	Currency              string `json:"currency"`
	CountryCurrencyDesc   string `json:"country_currency_desc"`
	ExchangeRate          string `json:"exchange_rate"`
	EffectiveDate         string `json:"effective_date"`
	SrcLineNbr            string `json:"src_line_nbr"`
	RecordFiscalYear      string `json:"record_fiscal_year"`
	RecordFiscalQuarter   string `json:"record_fiscal_quarter"`
	RecordCalendarYear    string `json:"record_calendar_year"`
	RecordCalendarQuarter string `json:"record_calendar_quarter"`
	RecordCalendarMonth   string `json:"record_calendar_month"`
	RecordCalendarDay     string `json:"record_calendar_day"`
}

type Meta struct {
	Count int `json:"count"`
}

type Response struct {
	Data []Record `json:"data"`
	Meta Meta     `json:"meta"`
}

func GetEntity(url string, headers map[string]string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar a requisição: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer a requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("requisição falhou com status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("erro ao decodificar a resposta JSON: %w", err)
	}

	return nil
}

func CapitalizeFirstLetter(text string) string {
	return cases.Title(language.Und, cases.Compact).String(text)
}

/*****
other funcs
******/
