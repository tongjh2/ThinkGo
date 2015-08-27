package ThinkGo

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	//"tongjh/tools"
)

type ControllerInterface interface {
	Init(w http.ResponseWriter, r *http.Request)
}

type Controller struct {
	Ctx   *Context
	Data  map[string]interface{}
	Funcs map[string]interface{}
}

func (c *Controller) Init(w http.ResponseWriter, r *http.Request) {
	c.Ctx = &Context{w, r}
	c.Data = make(map[string]interface{})
	c.Funcs = make(map[string]interface{})
}

func (c *Controller) Input() url.Values {
	if c.Ctx.Req.Form == nil {
		c.Ctx.Req.ParseForm()
	}
	return c.Ctx.Req.Form
}

//模板
func (c *Controller) Display(filename string) {
	t := template.New("name")
	if len(c.Funcs) > 0 {
		t = t.Funcs(c.Funcs)
	}
	bytes, err := ioutil.ReadFile(filename)
	t, err = t.Parse(string(bytes))
	fmt.Println(err)
	t.Execute(c.Ctx.Resp, c.Data)
}

//模板变量分配
func (this *Controller) Assign(key string, value interface{}) {
	this.Data[key] = value
	fmt.Println(this.Data)
}

//模板函数分配
func (this *Controller) Func(funcname string, f interface{}) {
	this.Funcs[funcname] = f
	fmt.Println(this.Funcs)
}

//输入输出
type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request
}

func (ctx *Context) WriteString(str string) {
	ctx.Resp.Write([]byte(str))
}

func (c *Controller) WriteString(str string) {
	c.Ctx.Resp.Write([]byte(str))
}
