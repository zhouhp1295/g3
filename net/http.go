// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package net

import (
	"github.com/gin-gonic/gin"
	"github.com/zhouhp1295/g3"
	"github.com/zhouhp1295/g3/auth"
	"github.com/zhouhp1295/g3/crud"
	"go.uber.org/zap"
	"net/http"
)

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

func Result(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func Success(ctx *gin.Context, msg string, data interface{}) {
	Result(ctx, http.StatusOK, msg, data)
}

func SuccessDefault(ctx *gin.Context) {
	Result(ctx, http.StatusOK, "Success", "")
}

func SuccessData(ctx *gin.Context, data interface{}) {
	Result(ctx, http.StatusOK, "Success", data)
}

func SuccessList(ctx *gin.Context, rows interface{}) {
	Result(ctx, http.StatusOK, "Success", map[string]interface{}{
		"rows": rows,
	})
}

func SuccessPage(ctx *gin.Context, rows interface{}, page crud.PageData) {
	Result(ctx, http.StatusOK, "Success", map[string]interface{}{
		"rows": rows,
		"page": page,
	})
}

func FailedServerError(ctx *gin.Context, msg string, data interface{}) {
	Result(ctx, http.StatusInternalServerError, msg, data)
}

func FailedBadRequest(ctx *gin.Context, msg string, data interface{}) {
	Result(ctx, http.StatusBadRequest, msg, data)
}

func FailedMessage(ctx *gin.Context, msg string) {
	Result(ctx, http.StatusBadRequest, msg, "")
}

func FailedNotFound(ctx *gin.Context) {
	Result(ctx, http.StatusBadRequest, "404 Not Found", "")
}

func ShouldBind(ctx *gin.Context, data interface{}) error {
	if http.MethodGet == ctx.Request.Method ||
		http.MethodDelete == ctx.Request.Method {
		return ctx.ShouldBindQuery(data)
	}
	return ctx.ShouldBindJSON(data)
}

func (baseApi *BaseApi) HandleGet(ctx *gin.Context) {
	params := IdParams{}
	err := ShouldBind(ctx, &params)
	if err != nil {
		g3.ZL().Error("parse params failed. please check")
		FailedMessage(ctx, "参数错误")
		return
	}
	if baseApi.Dao.CountByPk(params.Id) == 0 {
		g3.ZL().Error("record not exist. please check", zap.Int64("id", params.Id))
		FailedNotFound(ctx)
		return
	}
	m := baseApi.Dao.FindByPk(params.Id)

	baseApi.Dao.AfterGet(m)

	SuccessData(ctx, m)
}

func (baseApi *BaseApi) HandleInsert(ctx *gin.Context) {
	params := baseApi.Dao.GetModel().NewModel()
	err := ShouldBind(ctx, &params)
	if err != nil {
		g3.ZL().Error("parse params failed. please check")
		FailedMessage(ctx, "参数错误")
		return
	}
	if baseApi.Dao.CountByPk(params.GetId()) != 0 {
		g3.ZL().Error("duplicate primary key. please check")
		FailedMessage(ctx, "主键重复")
		return
	}
	if _ok, _msg := baseApi.Dao.BeforeInsert(params); !_ok {
		g3.ZL().Error("insert validate failed", zap.String("msg", _msg))
		FailedMessage(ctx, "操作失败:"+_msg)
		return
	}
	operator := ctx.GetInt64(auth.CtxJwtUid)
	if baseApi.Dao.Insert(params, operator) {
		baseApi.Dao.AfterInsert(params)
		SuccessData(ctx, params)
	} else {
		g3.ZL().Error("insert failed. please check", zap.Reflect("data", params))
		FailedMessage(ctx, "操作失败, 请稍后重试")
	}
}

func (baseApi *BaseApi) HandleUpdate(ctx *gin.Context) {
	params := baseApi.Dao.GetModel().NewModel()
	err := ShouldBind(ctx, &params)
	if err != nil {
		g3.ZL().Error("parse params failed. please check")
		FailedMessage(ctx, "参数错误")
		return
	}
	if baseApi.Dao.CountByPk(params.GetId()) == 0 {
		g3.ZL().Error("record not exist. please check", zap.Int64("id", params.GetId()))
		FailedNotFound(ctx)
		return
	}
	if _ok, _msg := baseApi.Dao.BeforeUpdate(params); !_ok {
		g3.ZL().Error("update validate failed", zap.String("msg", _msg))
		FailedMessage(ctx, "操作失败:"+_msg)
		return
	}
	operator := ctx.GetInt64(auth.CtxJwtUid)
	if baseApi.Dao.Update(params, operator) {
		baseApi.Dao.AfterUpdate(params)
		SuccessDefault(ctx)
	} else {
		g3.ZL().Error("update failed. please check", zap.Reflect("data", params))
		FailedMessage(ctx, "操作失败, 请稍后重试")
	}
}

func (baseApi *BaseApi) HandleUpdateStatus(ctx *gin.Context) {
	params := UpdateStatusParams{}
	err := ShouldBind(ctx, &params)
	if err != nil {
		g3.ZL().Error("parse params failed. please check")
		FailedMessage(ctx, "参数错误")
		return
	}
	if baseApi.Dao.CountByPk(params.Id) == 0 {
		g3.ZL().Error("record not exist. please check", zap.Int64("id", params.Id))
		FailedNotFound(ctx)
		return
	}
	if len(params.Status) == 0 {
		g3.ZL().Error("status is empty. please check")
		FailedMessage(ctx, "参数错误")
		return
	}
	operator := ctx.GetInt64(auth.CtxJwtUid)

	if baseApi.Dao.UpdateStatus(params.Id, params.Status, operator) {
		SuccessDefault(ctx)
	} else {
		g3.ZL().Error("update status failed. please check", zap.Reflect("data", params))
		FailedMessage(ctx, "操作失败, 请稍后重试")
	}
}

func (baseApi *BaseApi) HandleDelete(ctx *gin.Context) {
	params := IdParams{}
	err := ShouldBind(ctx, &params)
	if err != nil {
		g3.ZL().Error("parse params failed. please check")
		FailedMessage(ctx, "参数错误")
		return
	}
	if baseApi.Dao.CountByPk(params.Id) == 0 {
		g3.ZL().Error("record not exist. please check", zap.Int64("id", params.Id))
		FailedNotFound(ctx)
		return
	}

	m := baseApi.Dao.FindByPk(params.Id)

	if _ok, _msg := baseApi.Dao.BeforeDelete(m); !_ok {
		g3.ZL().Error("delete validate failed", zap.String("msg", _msg))
		FailedMessage(ctx, "操作失败:"+_msg)
		return
	}
	operator := ctx.GetInt64(auth.CtxJwtUid)
	if baseApi.Dao.Delete(m, operator) {
		baseApi.Dao.AfterDelete(m)
		SuccessDefault(ctx)
	} else {
		g3.ZL().Error("delete failed. please check", zap.Reflect("data", params))
		FailedMessage(ctx, "操作失败, 请稍后重试")
	}
}

func (baseApi *BaseApi) HandleRemove(ctx *gin.Context) {
	params := IdParams{}
	err := ShouldBind(ctx, &params)
	if err != nil {
		g3.ZL().Error("parse params failed. please check")
		FailedMessage(ctx, "参数错误")
		return
	}
	if baseApi.Dao.CountByPk(params.Id) == 0 {
		g3.ZL().Error("record not exist. please check", zap.Int64("id", params.Id))
		FailedNotFound(ctx)
		return
	}

	m := baseApi.Dao.FindByPk(params.Id)

	if _ok, _msg := baseApi.Dao.BeforeRemove(m); !_ok {
		g3.ZL().Error("remove validate failed", zap.String("msg", _msg))
		FailedMessage(ctx, "操作失败:"+_msg)
		return
	}
	operator := ctx.GetInt64(auth.CtxJwtUid)
	if baseApi.Dao.Remove(m, operator) {
		baseApi.Dao.AfterRemove(m)
		SuccessDefault(ctx)
	} else {
		g3.ZL().Error("remove failed. please check", zap.Reflect("data", params))
		FailedMessage(ctx, "操作失败, 请稍后重试")
	}
}

func (baseApi *BaseApi) HandleList(ctx *gin.Context) {
	modelParams := baseApi.Dao.GetModel().NewModel().(crud.ModelInterface)
	baseParams := new(crud.BaseQueryParams)
	_ = ShouldBind(ctx, modelParams)
	_ = ShouldBind(ctx, baseParams)
	rows := baseApi.Dao.FindList(modelParams, baseParams)
	SuccessList(ctx, rows)
}

func (baseApi *BaseApi) HandlePage(ctx *gin.Context) {
	modelParams := baseApi.Dao.GetModel().NewModel().(crud.ModelInterface)
	baseParams := new(crud.BaseQueryParams)
	_ = ShouldBind(ctx, modelParams)
	_ = ShouldBind(ctx, baseParams)
	rows, pageData := baseApi.Dao.FindPage(modelParams, baseParams)
	SuccessPage(ctx, rows, pageData)
}
