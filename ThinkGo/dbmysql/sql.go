package dbmysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//连接数据库
func conn() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/go?charset=utf8")
	if err != nil {
		return db, err
	}
	return db, err
}

//执行查询
func Query(sql string) (map[int]map[string]string, error) {
	db, err := conn()
	results := make(map[int]map[string]string)
	if err != nil {
		db.Close()
		return results, err
	}
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(sql)
		return results, err
	}

	cols, _ := rows.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}

	i := 0
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			return results, err
		}
		row := make(map[string]string)

		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row
		i++
	}

	rows.Close()
	db.Close()
	return results, nil
}

//执行增删改
func Exec(sqlString string) (sql.Result, error) {
	db, err := conn()
	if err != nil {
		return nil, err
	}
	stmt, err := db.Prepare(sqlString)

	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec()

	stmt.Close()
	db.Close()
	return result, err
}
