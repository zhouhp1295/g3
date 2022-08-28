// Copyright (c) 554949297@qq.com . 2022-2022. All rights reserved

package g3

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/zhouhp1295/g3/auth"
	"sync"
)

var (
	g3g         *Gin
	initGinOnce sync.Once
)

type RGroup struct {
	Group   *gin.RouterGroup
	path    string
	perms   *auth.Perm
	jwt     *auth.JwtAuth
	jwtOnce sync.Once
}

func (rg *RGroup) Path() string {
	return rg.path
}

func (rg *RGroup) NewJwt(secret string, expires int64) {
	rg.jwtOnce.Do(func() {
		rg.perms = auth.NewPerm()
		rg.jwt = auth.NewJwt(rg.path, rg.perms, secret, expires)
		rg.Group.Use(rg.jwt.Authentication)
	})
}

func (rg *RGroup) NewJwtToken(uid int64, roles string) (string, error) {
	if rg.jwt == nil {
		ZL().Error("create token failed ! jwt is nil.")
		return "", errors.New("jwt is nil")
	}
	return rg.jwt.Token(uid, roles)
}

func (rg *RGroup) Bind(method, router string, handler gin.HandlerFunc, perms ...string) {
	rg.Group.Handle(method, router, handler)
	if rg.perms != nil && rg.jwt != nil {
		if len(perms) == 0 {
			rg.MakeWhite(router)
		} else {
			rg.perms.AddRouterPerms(router, perms...)
		}
	}
}

func (rg *RGroup) MakeOpen(routers ...string) {
	if rg.jwt == nil {
		ZL().Error("make router to open failed ! jwt is nil.")
		return
	}
	rg.jwt.AddOpenRouters(routers...)
}

func (rg *RGroup) MakeWhite(routers ...string) {
	if rg.jwt == nil {
		ZL().Error("make router to white failed ! jwt is nil.")
		return
	}
	rg.jwt.AddWhiteRouters(routers...)
}

func (rg *RGroup) ClearRolePerm(role string) {
	if rg.perms == nil {
		ZL().Error("clear role perm failed ! perms is nil.")
		return
	}
	rg.perms.ClearRolePerm(role)
}

func (rg *RGroup) ClearAllRolesPerm() {
	if rg.perms == nil {
		ZL().Error("clear all roles perm failed ! perms is nil.")
		return
	}
	rg.perms.ClearAllRolesPerm()
}

func (rg *RGroup) AddRouterPerms(router string, perms ...string) {
	if rg.perms == nil {
		ZL().Error("add router perms failed ! perms is nil.")
		return
	}
	rg.perms.AddRouterPerms(router, perms...)
}

func (rg *RGroup) AddRolePerm(role string, perms ...string) {
	if rg.perms == nil {
		ZL().Error("add role perms failed ! perms is nil.")
		return
	}
	rg.perms.AddRolePerm(role, perms...)
}

type Gin struct {
	Engine *gin.Engine
	groups map[string]*RGroup
}

func (g *Gin) Html(relativePath string, handlers ...gin.HandlerFunc) {
	g.Engine.GET(relativePath, handlers...)
}

func (g *Gin) Group(path string) *RGroup {
	if group, exist := g.groups[path]; exist {
		return group
	}
	group := &RGroup{
		path:  path,
		Group: g.Engine.Group(path),
		perms: nil,
		jwt:   nil,
	}
	g.groups[path] = group
	return group
}

func GetGin() *Gin {
	if g3g == nil {
		panic("gin engine is nil")
	}
	return g3g
}

func SetGin(engine *gin.Engine) *Gin {
	initGinOnce.Do(func() {
		g3g = new(Gin)
		g3g.Engine = engine
		g3g.groups = make(map[string]*RGroup)
	})
	return g3g
}
