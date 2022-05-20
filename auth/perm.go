package auth

import (
	"github.com/zhouhp1295/g3/helpers"
	"strings"
	"sync"
)

const RootUser = "root"

const RootPerm = "*:*:*"

type Perm struct {
	rwMutex     *sync.RWMutex
	rolePerms   map[string][]string
	routerPerms map[string][]string
}

func NewPerm() *Perm {
	return &Perm{
		rwMutex:     new(sync.RWMutex),
		rolePerms:   make(map[string][]string),
		routerPerms: make(map[string][]string),
	}
}

func (p *Perm) AddRolePerm(role string, perms ...string) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	if len(perms) == 0 {
		return
	}
	if _, exist := p.rolePerms[role]; !exist {
		p.rolePerms[role] = make([]string, 0)
	}
	p.rolePerms[role] = append(p.rolePerms[role], perms...)
}

func (p *Perm) AddRouterPerms(router string, perms ...string) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	if len(perms) == 0 {
		return
	}
	if _, exist := p.routerPerms[router]; !exist {
		p.routerPerms[router] = make([]string, 0)
	}
	p.routerPerms[router] = append(p.routerPerms[router], perms...)
}

func (p *Perm) ReplaceRolePerm(role string, perms ...string) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.rolePerms[role] = perms
}

func (p *Perm) ClearRolePerm(role string) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.rolePerms[role] = make([]string, 0)
}

func (p *Perm) ClearAllRolesPerm() {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.rolePerms = make(map[string][]string)
}

func (p *Perm) CheckRolePerm(role string, perm string) bool {
	if role == RootUser {
		return true
	}
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	if rolePerms, ok2 := p.rolePerms[role]; ok2 {
		if helpers.IndexOf[string](rolePerms, perm) >= 0 {
			return true
		}
	}
	return false
}

func (p *Perm) CheckRolesPerm(roles string, perm string) bool {
	roleList := strings.Split(roles, ",")
	for _, role := range roleList {
		if p.CheckRolePerm(role, perm) {
			return true
		}
	}
	return false
}

func (p *Perm) CheckRoleRouter(role string, router string) bool {
	if role == RootUser {
		return true
	}
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	if routerPerms, ok := p.routerPerms[router]; ok {
		if rolePerms, ok2 := p.rolePerms[role]; ok2 {
			for _, routerPerm := range routerPerms {
				if helpers.IndexOf[string](rolePerms, routerPerm) >= 0 {
					return true
				}
			}
		}
		return false
	}
	return true
}

func (p *Perm) CheckRolesRouter(roles string, router string) bool {
	roleList := strings.Split(roles, ",")
	for _, role := range roleList {
		if p.CheckRoleRouter(role, router) {
			return true
		}
	}
	return false
}
