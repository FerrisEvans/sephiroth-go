package router

import (
	"sephiroth-go/router/biz"
	"sephiroth-go/router/sys"
)

type RouterGroup struct {
	Sys sys.RouterGroup
	Biz biz.RouterGroup
}
