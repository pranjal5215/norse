package backends

import (
	"github.com/goibibo/norse/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func incrFunc(iKey string) error { return nil }
func decrFunc(iKey string) error { return nil }

func init() {
	config.Configure("./test_config.json")
	Configure()
}

func Test_MySQL(t *testing.T) {
	mysqlClient, _ := GetMysqlClient(incrFunc, decrFunc, "mySQLConf")
	// "mysql" database name;
	_, err := mysqlClient.Select("select * from customers")
	assert.NoError(t, err)

	mysqlClient, _ = GetMysqlClient(incrFunc, decrFunc, "mySQLConf")
	// "mysql" database name;
	_, err = mysqlClient.Select("select * from customers")
	assert.NoError(t, err)
}
