package ThinkGo

import (
	//"ThinkGo/controller"
	"ThinkGo/tools"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

var controllerMap = make(map[string]interface{})

// func main() {
// 	http.Handle("/static/", http.FileServer(http.Dir("static")))
// 	http.HandleFunc("/", ThinkHandler)
// 	AutoRouter("admin", &controller.AdminController{})
// 	http.ListenAndServe(":8000", nil)
// }

func Run() {
	http.Handle("/static/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/", ThinkHandler)
	http.ListenAndServe(":8000", nil)
}

func AutoRouter(controllerName string, obj interface{}) {
	controllerMap[controllerName] = obj
}

func ThinkHandler(w http.ResponseWriter, r *http.Request) {
	pathinfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathinfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	c := parts[0]
	if !tools.In_mapSI(c, controllerMap) {
		w.Write([]byte("没有找到这个控制器"))
		return
	}

	controller_obj := controllerMap[c].(ControllerInterface)
	controller_obj.Init(w, r)

	controller := reflect.ValueOf(controller_obj)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
		if !method.IsValid() {
			w.Write([]byte(`没有找到 ` + action + ` 这个方法`))
			return
		}
	}
	argus := make([]reflect.Value, 0)
	fmt.Println(argus)
	method.Call(argus)
}
