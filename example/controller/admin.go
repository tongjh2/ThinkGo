package controller

import (
	"ThinkGo"
	"example/tools"
	"os"

	"fmt"
)

type AdminController struct {
	ThinkGo.Controller
}

func (this *AdminController) ShowAction() {
	file, err := os.Open("conf/app.conf")
	fmt.Println(file, err)
	data := make([]byte, 100)
	count, err := file.Read(data)
	fmt.Println(count, err, data, string(data))
	fmt.Println(os.Hostname())
	fmt.Println(os.Getpagesize())
	fmt.Println(os.Environ())

	fmt.Println(this.Input())
	fmt.Println(this.Ctx.Req.Form)
	this.Assign("data", "我是data变量")
	this.Func("test", tools.Test)
	this.Display("view/index.html")
}
