package controller

import (
	"ThinkGo/dbmysql"
	"ThinkGo/tools"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type AdminController struct {
}

func (this *AdminController) ShowAction(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/index.html")
	// data := map[string]string{"key1": "value1", "key2": "value2"}
	// checkError(err)
	// t = t.Funcs(template.FuncMap{"Test": tools.Test})
	// fmt.Println(tools.Test("aaaaa"))
	// t.Execute(w, data)

	// t, err := template.ParseFiles("views/index.html")
	// checkError(err)
	// data := map[string]string{"key1": "value1", "key2": "value2"}
	// t.Execute(w, data)

	tpl, _ := template.New("").Funcs(template.FuncMap{"Test": tools.Test}).ParseFiles("views/index.html")
	fmt.Println(tpl.Name())
	tpl.ExecuteTemplate(w, "index.html", nil)
	tpl.Execute(w, nil)

	w.Write([]byte("你好世界"))
}

func (this *AdminController) ListAction(w http.ResponseWriter, r *http.Request) {

	// db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/go?charset=utf8")
	// if err != nil {
	// 	panic(err.Error())
	// 	fmt.Println(err.Error())
	// }
	// defer db.Close()
	// stmt, err := db.Prepare()
	// defer stmt.Close()

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// result, err := dbmysql.Exec("INSERT INTO go_user(username, password) VALUES('aaa', 'bbbb')")

	// fmt.Println(result, err)
	// fmt.Println(result.LastInsertId())
	// fmt.Println(result.RowsAffected())

	db := dbmysql.NewM()
	data, err := db.Table("go_user").Select()
	fmt.Println(data, err)

	// data, err := dbmysql.Query("select * from go_user")
	// fmt.Println(data, err)

	// db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/go?charset=utf8")
	// if err != nil {
	// 	panic(err.Error())
	// 	fmt.Println(err.Error())
	// }
	// defer db.Close()
	// query, err := db.Query("select * from go_user")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer query.Close()

	// cols, _ := query.Columns()
	// values := make([][]byte, len(cols))
	// scans := make([]interface{}, len(cols))
	// for i := range values {
	// 	scans[i] = &values[i]
	// }

	// results := make(map[int]map[string]string)
	// i := 0
	// for query.Next() {
	// 	if err := query.Scan(scans...); err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println(scans)
	// 	row := make(map[string]string)

	// 	for k, v := range values {
	// 		key := cols[k]
	// 		row[key] = string(v)
	// 	}
	// 	results[i] = row
	// 	i++
	// }

	// for k, v := range results {
	// 	fmt.Println(k, v)
	// }

	// db.Close()
	w.Write([]byte("你好世界！"))
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
