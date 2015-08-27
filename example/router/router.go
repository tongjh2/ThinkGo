package router

import (
	"ThinkGo"
	"example/controller"
)

func init() {
	ThinkGo.AutoRouter("admin", &controller.AdminController{})
}
