package property

import (
	"database/sql"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"quick/manager/database"
	"quick/pkg/log"
	queue2 "quick/pkg/queue"
	"quick/pkg/tdengine"
	"runtime"
	"time"
)

var db *database.Database

var queue *queue2.Queue

func batch(q *queue2.Queue, workers int, taos *sql.DB) {
	for i := 0; i < workers; i++ {
		go func(q *queue2.Queue) {
			for {
				runtime.Gosched()
				msg, err := q.Dequeue()
				if err != nil {
					time.Sleep(200 * time.Millisecond)
					continue
				}
				sqlStr, err := tdengine.ToTaosBatchInsertSql(msg)
				if err != nil {
					log.Sugar.Errorf("cannot build sql with records: %v", err)
					continue
				}
				runtime.Gosched()
				//fmt.Println("sql:======" + sqlStr)
				_, err = taos.Exec(sqlStr)
				if err != nil {
					log.Sugar.Errorf("exec query error: %v, the sql command is:\n%s\n", err, sqlStr)
				}

			}
		}(q)
	}
}

var mqt mqtt.Client

func Publish(topic string, data interface{}) {
	payload, _ := json.Marshal(data)
	if token := mqt.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}
