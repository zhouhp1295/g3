// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package g3

import (
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

const (
	defaultAppName = "app"
	defaultAppId   = "default"
)

type Cfg struct {
	HomeDir string
	AppName string
	AppId   string
}

var (
	g3Cfg    *Cfg
	bootSync sync.Once
)

// Boot 初始化
func Boot(cfg *Cfg) {
	bootSync.Do(func() {
		if cfg == nil {
			g3Cfg = new(Cfg)
		} else {
			g3Cfg = cfg
		}
		if len(g3Cfg.AppName) == 0 {
			g3Cfg.AppName = defaultAppName
		}
		if len(g3Cfg.AppId) == 0 {
			g3Cfg.AppId = defaultAppId
		}
		if len(g3Cfg.HomeDir) == 0 {
			g3Cfg.HomeDir = filepath.Dir(AppPath())
		}
		defaultLogger = NewLogger(cfg.AppName, true)
	})
}

// HomeDir 工作目录
func HomeDir() string {
	if g3Cfg == nil {
		panic("please init g3 cfg")
	}
	return g3Cfg.HomeDir
}

func AppName() string {
	if g3Cfg == nil {
		panic("please init g3 cfg")
	}
	return g3Cfg.AppName
}

func AppId() string {
	if g3Cfg == nil {
		panic("please init g3 cfg")
	}
	return g3Cfg.AppId
}

// EnsureAbs prepends the HomeDir to the given path if it is not an absolute path.
func EnsureAbs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(HomeDir(), path)
}

func AssetPath(filename string) string {
	return filepath.Join(HomeDir(), filename)
}

var (
	appPath     string
	appPathOnce sync.Once
)

// AppPath app路径
func AppPath() string {
	appPathOnce.Do(func() {
		var err error
		appPath, err = exec.LookPath(os.Args[0])
		if err != nil {
			panic("look executable path: " + err.Error())
		}

		appPath, err = filepath.Abs(appPath)
		if err != nil {
			panic("get absolute executable path: " + err.Error())
		}
	})
	return appPath
}

var defaultLogger *zap.Logger

func ZL() *zap.Logger {
	return defaultLogger
}
