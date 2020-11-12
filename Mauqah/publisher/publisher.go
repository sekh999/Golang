package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	gopher_and_rabbit "rabbitmq"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}
func handleSimpleErr(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}

}

func main() {
	for {
		flag := true
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		input, _ := reader.ReadString('\n')
		log.Printf("Input: %s", input)
		conn, err := amqp.Dial(gopher_and_rabbit.Config.AMQPConnectionURL)
		handleError(err, "Can't connect to AMQP")
		defer conn.Close()

		amqpChannel, err := conn.Channel()
		handleError(err, "Can't create a amqpChannel")

		defer amqpChannel.Close()

		queue, err := amqpChannel.QueueDeclare("add", false, false, false, false, nil)
		handleError(err, "Could not declare `add` queue")

		rand.Seed(time.Now().UnixNano())
		addTaskReq := gopher_and_rabbit.AddTask{}
		errForm := json.Unmarshal([]byte(input), &addTaskReq)
		if errForm != nil {
			handleError(errForm, "Error in input Json format")
		}

		db, err := sql.Open("mysql", gopher_and_rabbit.MysqlConfig.MysqlURL)
		if err != nil {
			handleError(err, "Error in Connecting Mysql")
		}

		switch strings.ToUpper(addTaskReq.Operation) {
		case "CREATE":
			stmt, e := db.Prepare("insert into person(name, age) VALUES (?, ?)")
			if e != nil {
				handleError(e, "Error creating preparing stmt")
			}
			res, e := stmt.Exec(addTaskReq.Name, addTaskReq.Age)
			if e != nil {
				handleError(e, "Error inserting into Database")
			}
			id, e := res.LastInsertId()
			if e != nil {
				handleError(e, "Error getting Last insert id")
			}
			addTaskReq.RecordId = id
		case "UPDATE":
			if addTaskReq.RecordId == 0 {
				err2 := fmt.Errorf("Please specify recordId for updation")
				handleSimpleErr(err2, "")
				flag = false
			} else {
				stmt, e := db.Prepare("update person set name=? ,age=?  where id=?")
				if e != nil {
					handleError(e, "Error creating preparing stmt")
				}
				res, e := stmt.Exec(addTaskReq.Name, addTaskReq.Age, addTaskReq.RecordId)
				if e != nil {
					handleError(e, "Error updating record")
				}
				id, e := res.LastInsertId()
				addTaskReq.Response = " has been updated"
				addTaskReq.RecordId = id
			}
		case "DELETE":
			if addTaskReq.RecordId == 0 {
				err2 := fmt.Errorf("Please specify recordId for deletion")
				handleSimpleErr(err2, "")
				flag = false
			} else {
				stmt, e := db.Prepare("delete from person  where id=?")
				if e != nil {
					handleError(e, "Error creating preparing stmt")
				}
				stmt.Exec(addTaskReq.RecordId)

				addTaskReq.Response = "has been deleted"
			}
		case "GET":
			if addTaskReq.Name == "" {
				err2 := fmt.Errorf("Please specify Name for getting record")
				handleSimpleErr(err2, "")
				flag = false
			} else {
				rows, e := db.Query("select name,age,id from person  where name=?", addTaskReq.Name)
				if e != nil {
					handleError(e, "Error creating preparing stmt")
				}
				addTaskReq = gopher_and_rabbit.AddTask{}
				for rows.Next() {
					e = rows.Scan(&addTaskReq.Name, &addTaskReq.Age, &addTaskReq.RecordId)
					handleError(e, "Error fetching record")
					fmt.Println(addTaskReq)
				}
			}
		default:
			fmt.Println("Invalid Operation specified")
			flag = false
		}

		//addTask := gopher_and_rabbit.AddTask{input}
		if flag == true {
			body, err := json.Marshal(addTaskReq)
			if err != nil {
				handleError(err, "Error encoding JSON")
			}

			err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         body,
			})

			if err != nil {
				log.Fatalf("Error publishing message: %s", err)
			}
		}

		//log.Printf("AddTask: %s+%d", addTaskReq.Name, addTaskReq.Age)
	}

}
