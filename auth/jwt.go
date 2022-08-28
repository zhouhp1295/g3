// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/zhouhp1295/g3/helpers"
	"net/http"
	"strings"
	"time"
)

const (
	CtxJwtUid   = "CtxJwtUid"
	CtxJwtRoles = "CtxJwtRoles"
)

type jwtClaims struct {
	jwt.StandardClaims
	Uid         int64
	Roles       string
	ExpiredDate string
}

type JwtAuth struct {
	prefix       string
	perm         *Perm
	openApiList  []string //无需登录即可访问的接口
	whiteApiList []string //白名单, 登录后即可访问的接口
	expires      int64    //有效期,单位秒
	secret       string
}

func NewJwt(prefix string, perm *Perm, secret string, expires int64) *JwtAuth {
	return &JwtAuth{
		prefix:       prefix,
		perm:         perm,
		whiteApiList: make([]string, 0),
		expires:      expires,
		secret:       secret,
	}
}

func (jwtAuth *JwtAuth) AddWhiteRouters(routers ...string) {
	for _, router := range routers {
		if helpers.IndexOf[string](jwtAuth.whiteApiList, router) < 0 {
			jwtAuth.whiteApiList = append(jwtAuth.whiteApiList, router)
		}
	}
}

func (jwtAuth *JwtAuth) AddOpenRouters(routers ...string) {
	for _, router := range routers {
		if helpers.IndexOf[string](jwtAuth.openApiList, router) < 0 {
			jwtAuth.openApiList = append(jwtAuth.openApiList, router)
		}
	}
}

func (jwtAuth *JwtAuth) Token(uid int64, roles string) (string, error) {
	nowTime := time.Now()
	expiredTime := nowTime.Add(time.Duration(jwtAuth.expires) * time.Second)

	claims := jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
		Uid:         uid,
		Roles:       roles,
		ExpiredDate: helpers.FormatDefaultDate(expiredTime),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(jwtAuth.secret))

	return token, err
}

func (jwtAuth *JwtAuth) Parse(token string) (*jwtClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtAuth.secret), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwtClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func (jwtAuth *JwtAuth) Authentication(ctx *gin.Context) {
	_abort := func() {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "message": "Unauthorized"})
	}
	//开放接口校验
	router := strings.Replace(ctx.Request.URL.Path, jwtAuth.prefix, "", -1)
	if len(jwtAuth.openApiList) > 0 && helpers.IndexOf[string](jwtAuth.openApiList, router) >= 0 {
		ctx.Next()
		return
	}
	authToken := ctx.GetHeader("Authorization")
	if len(authToken) == 0 || !strings.HasPrefix(authToken, "Bearer ") {
		_abort()
		return
	}

	token := strings.Replace(authToken, "Bearer ", "", -1)

	claims, err := jwtAuth.Parse(token)

	if err != nil {
		_abort()
		return
	}

	if claims.Uid <= 0 {
		_abort()
		return
	}

	ctx.Set(CtxJwtUid, claims.Uid)
	ctx.Set(CtxJwtRoles, claims.Roles)

	// 白名单校验
	if len(jwtAuth.whiteApiList) > 0 && helpers.IndexOf[string](jwtAuth.whiteApiList, router) >= 0 {
		ctx.Next()
		return
	}
	// 权限校验
	if !jwtAuth.perm.CheckRolesRouter(claims.Roles, router) {
		_abort()
		return
	}

	ctx.Next()
}
