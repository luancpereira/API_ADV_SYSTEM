package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/luancpereira/APICheckout/apis/checkout/server/model/request"
	"github.com/luancpereira/APICheckout/apis/checkout/server/model/response"
	coreError "github.com/luancpereira/APICheckout/core/errors"
	"github.com/luancpereira/APICheckout/core/service"
)

type Checkout struct{}

/*****
funcs for posts
******/

// godoc
//
//	@Tags		Checkout Orders
//	@Produce	json
//	@Param		body	body		request.InsertTransaction	true	"Body JSON"
//	@Success	201		{object}	response.Created
//	@Failure	400		{object}	response.Exception
//	@Router		/api/checkout [post]
func (Checkout) InsertTransaction(ctx *gin.Context) {
	var req request.InsertTransaction
	err := GetBody(ctx, &req)
	if err != nil {
		return
	}

	ID, err := service.Checkout{}.CreateTransaction(req.Description, req.TransactionDate, req.TransactionValue)
	if err != nil {
		ResponseBadRequest(ctx, err)
		return
	}

	ResponseCreated(ctx, ID)
}

/*****
funcs for posts
******/

/*****
funcs for gets
******/

// godoc
//
//	@Tags		Checkout Orders
//	@Produce	json
//	@Param		transactionID	path		int64	true	"transactionID"
//	@Param		country			path		string	true	"country"
//	@Success	200				{object}	response.GetTransactionsByID
//	@Failure	400				{object}	response.Exception
//	@Router		/api/checkout/transactions/{transactionID}/country/{country} [get]
func (Checkout) GetByID(ctx *gin.Context) {
	transactionID, err := GetPathParamInt64(ctx, "transactionID", true)
	if err != nil {
		return
	}

	country, err := GetPathParamString(ctx, "country", true)
	if err != nil {
		return
	}

	model, err := service.Checkout{}.GetByID(transactionID, country)
	if err != nil {
		ResponseBadRequest(ctx, err)
		return
	}

	var res response.GetTransactionsByID
	err = copier.Copy(&res, model)
	if err != nil {
		ResponseBadRequest(ctx, err)
		return
	}

	ResponseOK(ctx, res)
}

// godoc
//
//	@Tags		Checkout Orders
//	@Produce	json
//	@Param		country					path		string	true	"country"
//	@Param		limit					query		int32	false	"limit min 1"	default(10)
//	@Param		offset					query		int32	false	"offset min 0"	default(0)
//	@Param		filter_transaction_date	query		string	true	"filter_transaction_date"
//	@Success	200						{object}	response.List{data=[]response.GetTransactions}
//	@Failure	400						{object}	response.Exception
//	@Router		/api/checkout/transactions/country/{country} [get]
func (Checkout) GetList(ctx *gin.Context) {
	country, err := GetPathParamString(ctx, "country", true)
	if err != nil {
		return
	}

	filters, _, limit, offset := GetQueryParam(ctx)
	filterTransactionDate := filters["transaction_date"]
	if filterTransactionDate == "" {
		ResponseBadRequest(ctx, coreError.New("error.transaction.date.required"))
		return
	}

	models, total, err := service.Checkout{}.GetList(filters, limit, offset, country)
	if err != nil {
		ResponseBadRequest(ctx, err)
		return
	}

	var res []response.GetTransactions

	for _, model := range models {
		res = append(res, response.GetTransactions{
			ID:                                      model.ID,
			Description:                             model.Description,
			TransactionDate:                         model.TransactionDate,
			TransactionValue:                        model.TransactionValue,
			ExchangeRate:                            model.ExchangeRate,
			TransactionValueConvertedToWishCurrency: model.TransactionValueConvertedToWishCurrency,
		})
	}

	ResponseListOk(ctx, res, total)
}

/*****
funcs for gets
******/

/*****
other funcs
******/

func GetBody(ctx *gin.Context, obj any) (err error) {
	err = ParseBody(ctx, obj)
	if err != nil {
		ResponseBadRequest(ctx, err)
		return
	}

	return
}

func ParseBody(ctx *gin.Context, obj any) (err error) {
	err = ctx.ShouldBindJSON(obj)
	if err != nil {
		err = coreError.New("error.request.body.invalid", err.Error())
		return
	}

	return
}

func GetQueryParam(ctx *gin.Context) (filters map[string]string, sorts map[string]string, limit, offset int64) {
	filters = make(map[string]string)
	sorts = make(map[string]string)
	limit = 10
	offset = 0

	for parameter, value := range ctx.Request.URL.Query() {
		if strings.HasPrefix(strings.ToLower(parameter), "sort_") {
			sorts[strings.ReplaceAll(parameter, "sort_", "")] = value[0]
		}

		if strings.HasPrefix(strings.ToLower(parameter), "filter_") {
			filters[strings.ReplaceAll(parameter, "filter_", "")] = value[0]
		}

		if parameter == "limit" {
			limitInt, _ := strconv.Atoi(value[0])
			limit = int64(limitInt)
		}

		if parameter == "offset" {
			offsetInt, _ := strconv.Atoi(value[0])
			offset = int64(offsetInt)
		}
	}

	return
}

func GetPathParamInt64(ctx *gin.Context, key string, required bool) (value int64, err error) {
	param, err := GetPathParamString(ctx, key, required)
	if err != nil {
		return
	}

	value, err = strconv.ParseInt(param, 10, 64)
	if err != nil {
		err = coreError.New("error.request.path.param.invalid", key)
		ResponseBadRequest(ctx, err)
	}

	return
}

func GetPathParamString(ctx *gin.Context, key string, required bool) (value string, err error) {
	value = ctx.Param(key)

	if required && len(strings.TrimSpace(value)) == 0 {
		err = coreError.New("error.request.path.param.invalid", key)
		ResponseBadRequest(ctx, err)
		return
	}

	return
}

func ResponseOK(ctx *gin.Context, bodyResponse any) {
	ctx.JSON(http.StatusOK, bodyResponse)
}

func ResponseListOk(ctx *gin.Context, bodyResponse any, total int64) {
	var list response.List

	list.Pagination = response.Pagination{Total: total}
	list.Data = bodyResponse

	ctx.JSON(http.StatusOK, list)
}

func ResponseCreated(ctx *gin.Context, ID int64) {
	bodyResponse := response.Created{ID: ID}

	ResponseCreatedBody(ctx, bodyResponse)
}

func ResponseCreatedBody(ctx *gin.Context, bodyResponse any) {
	ctx.JSON(http.StatusCreated, bodyResponse)
}

func ResponseBadRequest(ctx *gin.Context, err interface{}) {
	errOut := coreError.ConvertTo(err)

	ctx.AbortWithStatusJSON(http.StatusBadRequest, errOut)
}

/*****
other funcs
******/
