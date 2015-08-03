package backends

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goibibo/norse/config"
	"time"
)

//mysql wrapper struct
var config_map_mysql map[string]map[string]string

func configureMySql() {
	var err error
	config_map_mysql, err = config.LoadSqlConfig()
	if err != nil {
		panic("Error in loading config for mysql")

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
func GetMySql(incr, decr func(string) error, key string) *MySqlStruct {
	mysqldb := &sql.DB{}
	mysqldb.SetMaxOpenConns(10)
	mysqldb.SetMaxIdleConns(6)
	return &MySqlStruct{mysqldb, incr, decr, key}
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
func (m *MySqlStruct) Select(vertical, query string) ([]map[string]interface{}, error) {

	var err error
	url := getSQLUrl(vertical, config_map_mysql)
	m.DB, err = sql.Open("mysql", url)
	defer m.Close()
	if err != nil {
		return nil, err
	}
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
