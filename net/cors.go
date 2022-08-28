// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package net

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func DefaultCors() gin.HandlerFunc {
	return Cors("Authorization", "Access-Control-Allow-Origin")
}

func Cors(headers ...string) gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	cfg.AddAllowHeaders(headers...)
	cfg.AllowAllOrigins = true
	return cors.New(cfg)
}
