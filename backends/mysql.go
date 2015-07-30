package backends

import (
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/goibibo/norse/config"
)
//mysql wrapper struct
type MySqlStruct struct {
	*sql.DB
	incr func(string,int64)error
	decr func(string,int64)error
	key string
}

// Close mysql conn
func (m *MySqlStruct) Close(){
	_ = m.DB.Close()
}
// Your instance type for mysql
func GetMySql(incr,decr func(string,int64)error,key string) (*MySqlStruct) {
	return &MySqlStruct{&sql.DB{},incr,decr,key}
}

func (m *MySqlStruct)Execute(query string ) (*sql.Rows,error) {
	m.incr(m.key,1)
	defer m.decr(m.key,1)
	return m.DB.Query(query)
}
func getSQLUrl(vertical string ,config_map map[string] map[string] string) string{
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",config_map[vertical]["username"],config_map[vertical]["password"],config_map[vertical]["host"],config_map[vertical]["port"],config_map[vertical]["database"])

}
func incr(s string,i int64)error{
        return nil
}
func decr(s string,i int64)error{
        return nil
}
func (m *MySqlStruct) Select(vertical,query string) ([]map[string]interface{}, error) {

	var err error
	config_map,err:= config.LoadSqlConfig()
	 if err != nil {
                return nil, err
        }

	url :=getSQLUrl(vertical,config_map)
	m.DB, err = sql.Open("mysql",url)
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

