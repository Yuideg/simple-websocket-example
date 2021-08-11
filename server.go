package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/pgxpool"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net/http"
	"os"
)

func main() {
	DATABASE_URL := "postgres://postgres:yideg2378@localhost:5432/chat"
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	fmt.Println("database chat connected successfuly!")

	r := gin.Default()
	m := melody.New()
	path := []string{"/a", "/b", "/c", "/d"}
	for i := 0; i < len(path); i++ {
		r.GET("/"+path[i], func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "index.html")
		})
	}

	for i := 0; i < len(path); i++ {
		r.GET("/ws"+path[i], func(c *gin.Context) {
			err := m.HandleRequest(c.Writer, c.Request)
			if err != nil {
				log.Fatal(err)
				return
			}
		})
	}

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println("msg :", string(msg))
		var messages = GetItems(pool, ctx)
		fmt.Println("message from db ", messages)

		var mess interface{}
		err := json.Unmarshal(msg, &mess)
		fmt.Println("mess interface ", mess)
		if err != nil {
			log.Fatal("error ", err)
			return
		}

		info := Message{}
		err = mapstructure.Decode(mess, &info)
		fmt.Println("error decode ", err)
		if err != nil {
			log.Fatal("error decode ", err)
		}
		messages = append(messages, info)
		fmt.Println("messages ", messages)
		for i := 0; i < len(messages); i++ {
			mm := messages[i]
			fmt.Println("Message struct ", mm)
			err := m.Broadcast([]byte(mm.message))
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		NewItem(pool, ctx, mess)

	})

	err = r.Run(":5000")
	if err != nil {
		log.Fatal(err)
	}
}

type Message struct {
	message string `json:"message"`
}

//Add new Item the store
func NewItem(pgx *pgxpool.Pool, ctx context.Context, new_item interface{}) {
	//aa,ok:=new_item.(Message)
	mess := Message{}
	mapstructure.Decode(new_item, &mess)
	fmt.Println("new message ", mess)
	ProfileInsert := `INSERT INTO temp(message) VALUES($1);`
	_, err := pgx.Exec(ctx, ProfileInsert, mess.message)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("data inserted successfully")

}

//Update the existing  tem
func UpdateItem(pgx *pgxpool.Pool, ctx context.Context, item interface{}) {
	ProfileUpdate := `UPDATE temp SET message =$1;`
	//var m Message
	//m:=Message{}
	m := Message{}
	mapstructure.Decode(item, &m)
	_, err := pgx.Exec(ctx, ProfileUpdate, m.message)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("data updated successfully")

}

//List all Items from the store
func GetItems(pgx *pgxpool.Pool, ctx context.Context) []Message {
	ProfileQuery := `SELECT * FROM temp;`
	rows, err := pgx.Query(ctx, ProfileQuery)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	var information []Message
	for rows.Next() {
		info := Message{}
		err = rows.Scan(&info.message)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		information = append(information, info)
	}
	defer rows.Close()
	fmt.Println(information)
	fmt.Println("Information ", information)
	return information
}
