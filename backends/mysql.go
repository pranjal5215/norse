package main

import (
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type MySqlStruct struct {
	*sql.DB
}

// Close redis conn
func (m *MySqlStruct) Close(){
	_ = m.DB.Close()
}
// Your instance type for redis
func GetMySql() (*MySqlStruct) {
	return &MySqlStruct{}
}
type Closure struct {
	input [] interface{}
	Do func()
	output interface{}
	
}
func (m *MySqlStruct)Execute(incr,decr func(),query string ) (*sql.Rows,error) {
	incr()
	defer decr()
	return m.DB.Query(query)
}
func getSQLUrl(level1,level2 string ,m map[string]interface{})string{
	return ""
}
func (m *MySqlStruct) Select(level1,level2,query string) ([]map[string]interface{}, error) {


	var err error
	m.DB, err = sql.Open("mysql","")
	
	if err != nil {
		return nil, err
	}
	rows, err := m.Execute(nil,nil,query)
	if err != nil {
		fmt.Println(err)
		panic("From MySQL: Error in executing select query")
	}
	columns, _ := rows.Columns()

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		rows.Scan(scanArgs...)
		record := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				switch col.(type) {
				default:
					panic("From MySQL: Unknown Type in type switching")
				case bool:
					record[columns[i]] = col.(bool)
				case int:
					record[columns[i]] = col.(int)
				case int64:
					record[columns[i]] = col.(int64)
				case float64:

					record[columns[i]] = col.(float64)
				case string:
					record[columns[i]] = col.(string)
				case []byte: // -- all cases go HERE!
					record[columns[i]] = string(col.([]byte))
				case time.Time:
					record[columns[i]] = col.(string)
				}
			}
		}
		records = append(records, record)
	}
	return records, nil
}
// How to use mysql,
func main(){
	rStruct := GetMySql()
	query:=""
	value, _ := rStruct.Select("goibibo", query,"")
	fmt.Println(value)
}
