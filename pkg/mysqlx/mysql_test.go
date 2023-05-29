package mysqlx

import (
	"log"
	"testing"

	"github.com/MisterChing/go-lib/utils/debugutil"
)

func TestMysql(t *testing.T) {

	dsn := "user:password@tcp(host:3306)/xxx?charset=utf8mb4&parseTime=true&timeout=3s&loc=Local"
	db, err := NewDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	debugutil.DebugPrintV2("111", err, db)

}
