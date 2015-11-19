package backends

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/goibibo/norse/config"
)

//mysql wrapper struct
var config_map_mysql map[string]map[string]string
var mysqlstructmap map[string]*MySqlStruct

func init() {
	mysqlstructmap = make(map[string]*MySqlStruct)
}

// mysqlstructmap is a map of pools
func configureMySql() {
	var err error
	config_map_mysql, err = config.LoadSqlConfig()
	if err != nil {
		panic("Error in loading config for mysql")

	}
	for vertical, _ := range config_map_mysql {
		url := getSQLUrl(vertical, config_map_mysql)
		mysqldb, err := sql.Open("mysql", url)
		if err != nil {
			panic("error in creating pool ")
		}
		mysqldb.SetMaxOpenConns(10)
		mysqldb.SetMaxIdleConns(6)
		mysqlstructmap[vertical] = &MySqlStruct{mysqldb, nil, nil, ""}
	}
}

type MySqlStruct struct {
	*sql.DB
	incr func(string) error
	decr func(string) error
	key  string
}

// Close mysql conn
func (m *MySqlStruct) Close() {
	_ = m.DB.Close()
}

// Your instance type for mysql
func GetMysqlClient(incr, decr func(string) error, vertical string) (*MySqlStruct, error) {
	sqlstruct, present := mysqlstructmap[vertical]
	if present {
		sqlstruct.incr = incr
		sqlstruct.decr = decr
		sqlstruct.key = vertical

		return sqlstruct, nil
	} else {
		panic("vertical not present in config")
	}
}

func (m *MySqlStruct) Execute(query string) (*sql.Rows, error) {
	m.incr(m.key)
	defer m.decr(m.key)
	return m.DB.Query(query)
}
func getSQLUrl(vertical string, config_map_mysql map[string]map[string]string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config_map_mysql[vertical]["username"], config_map_mysql[vertical]["password"], config_map_mysql[vertical]["host"], config_map_mysql[vertical]["port"], config_map_mysql[vertical]["database"])

}
func incr(s string, i int64) error {
	return nil
}
func decr(s string, i int64) error {
	return nil
}
func (m *MySqlStruct) Select(query string) ([]map[string]interface{}, error) {

	var err error
	rows, err := m.Execute(query)
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
				case []byte:
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
