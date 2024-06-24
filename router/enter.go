package router

import (
	"sephiroth-go/router/auth"
	"sephiroth-go/router/storage"
	"sephiroth-go/router/sys"
)

type RouterGroup struct {
	Auth    auth.RouterGroup
	Storage storage.RouterGroup
	Sys     sys.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
