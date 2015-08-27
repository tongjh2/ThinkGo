package dbmysql

import (
	"ThinkGo/tools"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var prefix = ""
var SQL string

func NewM() *DbMysql {
	return new(DbMysql)
}

func NewWhere() map[string]interface{} {
	return make(map[string]interface{})
}

type DbMysql struct {
	tableName string
	alias     string
	field     string
	join      string
	where     string
	order     string
	limit     string
	set       string
}

func (this *DbMysql) Alias(str string) *DbMysql {
	this.alias = str
	return this
}

func (this *DbMysql) Field(str string) *DbMysql {
	this.field = str
	return this
}

func (this *DbMysql) Table(str string) *DbMysql {
	this.tableName = prefix + str
	return this
}

func (this *DbMysql) Join(str string) *DbMysql {
	str = strings.Replace(str, "prefix_", prefix, -1)
	this.join += str + `
	`
	return this
}

func (this *DbMysql) Where(val interface{}) *DbMysql {
	switch value := val.(type) {
	case string:
		this.where += value
	case map[string]interface{}:
		var where string
		for k, v := range value {

			ks := strings.Split(k, "__")

			if len(ks) == 2 {

				arr := strings.Split(ks[0], "_")
				lianjie := "AND" //条件连接符
				operator := "eq" //关系运算符
				if len(arr) == 2 {
					if arr[0] == "or" {
						lianjie = "OR"
					}
					operator = arr[1]
				} else {
					operator = arr[0]
				}
				whereStr := "("
				if operator == "eq" {
					whereStr += ks[1] + "='" + Guolv(ToString(v)) + "'"
				} else if operator == "neq" {
					whereStr += ks[1] + " <> '" + Guolv(ToString(v)) + "'"
				} else if operator == "not" {
					whereStr += "not " + ks[1] + "='" + Guolv(ToString(v)) + "'"
				} else if operator == "lt" {
					whereStr += ks[1] + " < '" + Guolv(ToString(v)) + "'"
				} else if operator == "lte" {
					whereStr += ks[1] + " <= '" + Guolv(ToString(v)) + "'"
				} else if operator == "gt" {
					whereStr += ks[1] + " > '" + Guolv(ToString(v)) + "'"
				} else if operator == "gte" {
					whereStr += ks[1] + " >= '" + Guolv(ToString(v)) + "'"
				} else if operator == "notnull" {
					whereStr += ks[1] + " <> ''"
				} else if operator == "like" {
					whereStr += ks[1] + " like " + Guolv(ToString(v))
				} else if operator == "in" {
					whereStr += ks[1] + " in(" + strings.Replace(Guolv(ToString(v)), `\,`, ",", -1) + ")"
				} else if operator == "between" {
					arr1 := strings.Split(v.(string), ",")
					whereStr += ks[1] + " between '" + Guolv(ToString(arr1[0])) + "' and '" + Guolv(ToString(arr1[1])) + "'"
				}
				whereStr += ")"
				where += lianjie + " " + whereStr + " "

			} else if len(ks) == 1 {
				if ks[0] == "sql" {
					where += v.(string)
				} else {
					where += " AND (" + ks[0] + "='" + Guolv(ToString(v)) + "') "
				}
			}
		}
		this.where += where
	default:
		this.where += ""

	}
	return this
}

func (this *DbMysql) Order(str string) *DbMysql {
	this.order += str
	return this
}

func (this *DbMysql) Limit(str string) *DbMysql {
	this.limit += str
	return this
}

func (this *DbMysql) Ssql() string {
	sql := `select 
`
	if this.field == "" {
		this.field = ` * `
	}
	sql += this.field
	sql += ` 
from ` + this.tableName
	if this.alias != "" {
		sql += ` AS ` + this.alias
	}
	if this.join != "" {
		sql += ` 
	` + this.join
	}
	if this.where != "" {
		sql += `
where 1 ` + this.where
	}
	if this.order != "" {
		sql += ` 
	order by ` + this.order
	}

	return sql
}

func (this *DbMysql) Select() (data map[int]map[string]string, err error) {
	sql := this.Ssql()
	if this.limit == "" {
		sql += ` 
limit 0,100`
	} else {
		sql += ` 
limit ` + this.limit
	}
	SQL = sql
	fmt.Println(SQL)
	data, err = Query(sql)
	return data, err
}

func (this *DbMysql) Count() int64 {
	sql := this.Ssql() + ` 
limit 0,1`
	SQL = sql

	result, err := Query(SQL)
	if err != nil || len(result) == 0 {
		return 0
	}
	return int64(len(result))
}

func (this *DbMysql) One() (map[string]string, error) {
	sql := this.Ssql() + ` 
limit 0,1`
	SQL = sql

	var result map[string]string
	data, err := Query(SQL)
	if len(data) > 0 {
		result = data[0]
	}
	return result, err
}

func (this *DbMysql) Set(data map[string]interface{}) *DbMysql {
	var set string
	for k, v := range data {
		switch val := v.(type) {
		case string:
			set += `,` + k + `='` + Guolv(val) + `'`
		case int64:
			set += `,` + k + `=` + ToString(val)
		case int:
			set += `,` + k + `=` + ToString(val)
		}
	}
	if set != "" {
		this.set = tools.Substr(set, 1, len(set))
	}
	return this
}

func (this *DbMysql) Update(data map[string]interface{}) (sql.Result, error) {
	if len(data) > 0 {
		this.Set(data)
	}
	if this.where == "" || this.set == "" {
		return nil, errors.New("没有更新的参数")
	}
	SQL = `update ` + this.tableName + ` ` + this.alias + `
set ` + this.set + `
where 1 ` + this.where

	result, err := Exec(SQL)
	return result, err
}

func (this *DbMysql) Add(data map[string]interface{}) (sql.Result, error) {
	if len(data) > 0 {
		return nil, errors.New("传入的数据不是合法的")
	}

	var keys []string
	var values []string
	for k, v := range data {
		keys = append(keys, k)
		values = append(values, `'`+Guolv(ToString(v))+`'`)
	}
	key := strings.Join(keys, ",")
	value := strings.Join(values, ",")
	SQL = `insert into ` + this.tableName + `
(` + key + `)value
(` + value + `)`

	result, err := Exec(SQL)
	return result, err
}

func (this *DbMysql) Addall(data []map[string]interface{}) (sql.Result, error) {
	if len(data) > 0 {
		return nil, errors.New("传入的数据不是合法的")
	}

	var keys []string
	var vals []string
	for k, _ := range data[0] {
		keys = append(keys, k)
	}

	for _, b := range data {
		var values []string
		for _, v := range keys {
			values = append(values, `'`+Guolv(ToString(b[v]))+`'`)
		}
		vals1 := "(" + strings.Join(values, ",") + ")"
		vals = append(vals, vals1)
	}
	key := strings.Join(keys, ",")
	value := strings.Join(vals, ",")
	SQL = `insert into ` + this.tableName + `
(` + key + `)value
` + value

	result, err := Exec(SQL)
	return result, err
}

func (this *DbMysql) Delete(args ...int64) (sql.Result, error) {

	if len(args) > 0 {
		var ids []string
		for _, id := range args {
			ids = append(ids, ToString(id))
		}
		idstr := strings.Join(ids, ",")
		SQL = `delete from ` + this.tableName + ` where id in(` + ToString(idstr) + `)`
		result, err := Exec(SQL)
		return result, err
	}

	if this.where == "" {
		return nil, errors.New("必须指定条件才允许删除")
	}

	SQL = `delete from ` + this.tableName + ` where 1` + this.where
	result, err := Exec(SQL)
	return result, err
}

//sql字符串过滤
func Guolv(str string) string {
	str = strings.ToLower(str)
	str = strings.Replace(str, `'`, `\'`, -1)
	str = strings.Replace(str, `"`, `\"`, -1)
	str = strings.Replace(str, "`", "\\`", -1)
	str = strings.Replace(str, `;`, `\;`, -1)
	str = strings.Replace(str, `,`, `\,`, -1)
	str = strings.Replace(str, `:`, `\:`, -1)
	str = strings.Replace(str, `:`, `\:`, -1)
	str = strings.Replace(str, `#`, `\#`, -1)
	str = strings.Replace(str, `*`, `\*`, -1)
	str = strings.Replace(str, `/`, `\/`, -1)
	str = strings.Replace(str, `[`, `\[`, -1)
	str = strings.Replace(str, `]`, `\]`, -1)
	str = strings.Replace(str, `=`, `\=`, -1)
	str = strings.Replace(str, `or`, `\or`, -1)
	str = strings.Replace(str, `(`, `\(`, -1)
	return str
}

//解析一个where的map
func Get_where(where map[string]interface{}) string {
	str := ""
	if len(where) < 1 {
		return str
	}
	for k, v := range where {
		ks := strings.Split(k, "__")

		if len(ks) == 2 {

			arr := strings.Split(ks[0], "_")
			lianjie := "AND" //条件连接符
			operator := "eq" //关系运算符
			if len(arr) == 2 {
				if arr[0] == "or" {
					lianjie = "OR"
				}
				operator = arr[1]
			} else {
				operator = arr[0]
			}
			whereStr := "("
			if operator == "eq" {
				whereStr += ks[1] + "='" + Guolv(ToString(v)) + "'"
			} else if operator == "neq" {
				whereStr += ks[1] + " <> '" + Guolv(ToString(v)) + "'"
			} else if operator == "not" {
				whereStr += "not " + ks[1] + "='" + Guolv(ToString(v)) + "'"
			} else if operator == "lt" {
				whereStr += ks[1] + " < '" + Guolv(ToString(v)) + "'"
			} else if operator == "lte" {
				whereStr += ks[1] + " <= '" + Guolv(ToString(v)) + "'"
			} else if operator == "gt" {
				whereStr += ks[1] + " > '" + Guolv(ToString(v)) + "'"
			} else if operator == "gte" {
				whereStr += ks[1] + " >= '" + Guolv(ToString(v)) + "'"
			} else if operator == "notnull" {
				whereStr += ks[1] + " <> ''"
			} else if operator == "like" {
				whereStr += ks[1] + " like " + Guolv(ToString(v))
			} else if operator == "in" {
				whereStr += ks[1] + " in(" + strings.Replace(Guolv(ToString(v)), `\,`, ",", -1) + ")"
			} else if operator == "between" {
				arr1 := strings.Split(v.(string), ",")
				whereStr += ks[1] + " between '" + Guolv(ToString(arr1[0])) + "' and '" + Guolv(ToString(arr1[1])) + "'"
			}
			whereStr += ")"
			str += lianjie + " " + whereStr

		} else if len(ks) == 1 {
			if ks[0] == "sql" {
				str += v.(string)
			} else {
				str += " AND (" + ks[0] + "='" + Guolv(ToString(v)) + "')"
			}
		}
	}
	str = strings.TrimLeft(str, " AND ")
	str = strings.TrimLeft(str, " OR ")
	return ` (` + str + `) `
}

func ToString(args ...interface{}) string {
	result := ""
	for _, arg := range args {
		switch val := arg.(type) {
		case int:
			result += strconv.Itoa(val)
		case int8:
			result += strconv.Itoa(int(val))
		case int64:
			result += strconv.Itoa(int(val))
		case string:
			result += val
		case float64:
			result += strconv.FormatFloat(val, 'f', -1, 64)
		case float32:
			result += strconv.FormatFloat(float64(val), 'f', -1, 64)
		case time.Time:
			result += val.Format("2006-01-02 15:04:05")
		}
	}
	return result
}
