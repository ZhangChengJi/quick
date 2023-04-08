package tdengine

import (
	"database/sql"
	"fmt"
	"quick/conf"
	"quick/pkg/log"

	_ "github.com/taosdata/driver-go/v3/taosSql"
	"strconv"
)

func New(conf *conf.Tdengine) (*sql.DB, error) {
	url := "root:taosdata@/tcp(" + conf.Host + ":" + strconv.Itoa(conf.Port) + ")/"

	//open connect to taos server
	db, err := sql.Open("taosSql", url)
	if err != nil {
		log.Sugar.Fatalf("Open database error: %s\n", err)
	}
	err = db.Ping()
	fmt.Println("🪅Tdengine Connection Successful........")
	if err != nil {
		panic("tdengine conn failed")
		return nil, err
	}
	return db, nil
}
