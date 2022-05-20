package net

import (
	"github.com/gin-gonic/gin"
	"github.com/zhouhp1295/g3/crud"
	"net/http"
)

type Method string

type ApiInterface interface {
	Result(ctx *gin.Context, code int, msg string, data interface{})
	HandleGet(ctx *gin.Context)
	HandleInsert(ctx *gin.Context)
	HandleUpdate(ctx *gin.Context)
	HandleUpdateStatus(ctx *gin.Context)
	HandleDelete(ctx *gin.Context)
	HandleRemove(ctx *gin.Context)
	HandleList(ctx *gin.Context)
	HandlePage(ctx *gin.Context)
}

type IdParams struct {
	Id int64 `json:"id" form:"id"`
}

type UpdateStatusParams struct {
	IdParams
	Status string `json:"status" form:"status"`
}

type BaseApi struct {
	Dao crud.DAOInterface
}

func (baseApi *BaseApi) Result(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(code, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func (baseApi *BaseApi) Success(ctx *gin.Context, msg string, data interface{}) {
	baseApi.Result(ctx, http.StatusOK, msg, data)
}

func (baseApi *BaseApi) SuccessDefault(ctx *gin.Context) {
	baseApi.Result(ctx, http.StatusOK, "Success", "")
}

func (baseApi *BaseApi) SuccessData(ctx *gin.Context, data interface{}) {
	baseApi.Result(ctx, http.StatusOK, "Success", data)
}

func (baseApi *BaseApi) SuccessList(ctx *gin.Context, rows interface{}) {
	baseApi.Result(ctx, http.StatusOK, "Success", map[string]interface{}{
		"rows": rows,
	})
}

func (baseApi *BaseApi) SuccessPage(ctx *gin.Context, rows interface{}, page crud.PageData) {
	baseApi.Result(ctx, http.StatusOK, "Success", map[string]interface{}{
		"rows": rows,
		"page": page,
	})
}

func (baseApi *BaseApi) FailedServerError(ctx *gin.Context, msg string, data interface{}) {
	baseApi.Result(ctx, http.StatusInternalServerError, msg, data)
}

func (baseApi *BaseApi) FailedBadRequest(ctx *gin.Context, msg string, data interface{}) {
	baseApi.Result(ctx, http.StatusBadRequest, msg, data)
}

func (baseApi *BaseApi) FailedMessage(ctx *gin.Context, msg string) {
	baseApi.Result(ctx, http.StatusBadRequest, msg, "")
}

func (baseApi *BaseApi) FailedNotFound(ctx *gin.Context) {
	baseApi.Result(ctx, http.StatusBadRequest, "404 Not Found", "")
}

func (baseApi *BaseApi) HandleGet(ctx *gin.Context) {
	params := IdParams{}
	var err error
	switch Method(ctx.Request.Method) {
	case http.MethodGet:
		if err = ctx.ShouldBindQuery(&params); err != nil {
			baseApi.FailedMessage(ctx, "参数错误")
			return
		}
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		if err = ctx.ShouldBindJSON(&params); err != nil {
			baseApi.FailedMessage(ctx, "参数错误")
			return
		}
	default:
		baseApi.FailedMessage(ctx, "参数错误")
		return
	}

	if baseApi.Dao.CountByPk(params.Id) == 0 {
		baseApi.FailedNotFound(ctx)
		return
	}

	m := baseApi.Dao.FindByPk(params.Id)

	baseApi.SuccessData(ctx, m)
}

func (baseApi *BaseApi) HandleInsert(ctx *gin.Context) {

}

func (baseApi *BaseApi) HandleUpdate(ctx *gin.Context) {

}

func (baseApi *BaseApi) HandleUpdateStatus(ctx *gin.Context) {

}

func (baseApi *BaseApi) HandleDelete(ctx *gin.Context) {

}

func (baseApi *BaseApi) HandleRemove(ctx *gin.Context) {

}

func (baseApi *BaseApi) HandleList(ctx *gin.Context) {

}

func (baseApi *BaseApi) HandlePage(ctx *gin.Context) {
	modelParams := baseApi.Dao.GetModel().NewModel().(crud.ModelInterface)
	baseParams := new(crud.BaseQueryParams)
	switch Method(ctx.Request.Method) {
	case http.MethodGet:
		_ = ctx.ShouldBindQuery(modelParams)
		_ = ctx.ShouldBindQuery(baseParams)
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		_ = ctx.ShouldBindJSON(modelParams)
		_ = ctx.ShouldBindJSON(baseParams)
	default:
		baseApi.FailedMessage(ctx, "Method错误")
		return
	}
	rows, pageData := baseApi.Dao.FindPage(modelParams, baseParams)
	baseApi.SuccessPage(ctx, rows, pageData)
}
