package main

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/jackc/pgx/pgxpool"
//	"gopkg.in/olahol/melody.v1"
//	"log"
//	"os"
//)
//
//var clients = make(map[string]interface{})
//var sessions []*melody.Session
//
//
//
//
//
//
//
//
//
//func main() {
//	DATABASE_URL := "postgres://postgres:yideg2378@localhost:5432/chat"
//	ctx := context.Background()
//	pool, err := pgxpool.Connect(ctx, DATABASE_URL)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//	defer pool.Close()
//	fmt.Println("database chat connected successfuly!")
//
//	r := gin.Default()
//	m := melody.New()
//	path := []string{"/a", "/b", "/c", "/d"}
//	//for i := 0; i < len(path); i++ {
//	//	r.GET("/"+path[i], func(c *gin.Context) {
//	//		http.ServeFile(c.Writer, c.Request, "index.html")
//	//	})
//	//}
//	//var mel  *melody.Melody
//	//router:=pkg.NewName()
//	for i := 0; i < len(path); i++ {
//		r.GET("/ws"+path[i], func(c *gin.Context) {
//			err := m.HandleRequest(c.Writer, c.Request)
//			fmt.Println("error",err)
//			if err != nil {
//				log.Fatal("Unable to convert the http connection to websocket")
//				return
//			}
//		})
//
//	}
//
//	m.HandleDisconnect(func(s *melody.Session) {
//		for index, ss := range sessions {
//			if s == ss {
//				sessions = append(sessions[:index], sessions[(index+1):]...)
//			}
//		}
//		client := s.Request.URL.Path
//		fmt.Println("client ", clients[client], " is going to disconnected from the group")
//		delete(clients, client)
//		fmt.Println("clients ", clients, "no session ", len(sessions))
//	})
//	m.HandleConnect(func(s *melody.Session) {
//		sessions = append(sessions, s)
//		fmt.Println("New seesion is Added to Seesion collections", len(sessions))
//		fmt.Println("connecting ...")
//		host := s.Request.Host + s.Request.URL.Path
//		key := s.Request.URL.Path
//		clients[key] = host
//		s.Set(s.Request.URL.Path, host)
//		//fmt.Println("host address url",s.Keys[s.Request.URL.Host])
//		fmt.Println(host, " connected", "clients ", clients, "no session ", len(sessions))
//
//	})
//
//	m.HandleMessage(func(s *melody.Session, msg []byte) {
//		//x_url, ok := s.Get(s.Request.URL.Path)
//		var messages = GetItems(pool, ctx)
//
//		var mess interface{}
//		err := json.Unmarshal(msg, &mess)
//		if err != nil {
//			log.Fatal("error ", err)
//			return
//		}
//
//		info := Message{}
//		mes := mess.(map[string]interface{})
//		for _, value := range mes {
//			info.message = value.(string)
//		}
//		messages = append(messages, info)
//		var ses []*melody.Session
//		current_connected:=s.Keys[s.Request.URL.Path]
//
//		for _, se := range sessions {
//			find:=se.Keys[s.Request.URL.Path]
//			if  find==current_connected {
//				ses=append(ses,se)
//			}
//		}
//
//		for i := 0; i < len(messages); i++ {
//			mm := messages[i]
//			err := m.BroadcastMultiple([]byte(mm.message),ses)
//			if err != nil {
//				log.Fatal(err)
//				return
//			}
//		}
//		NewItem(pool, ctx, mes)
//	})
//	err = r.Run(":5000")
//	if err != nil {
//		log.Fatal(err)
//	}
//}

//type Message struct {
//	message string `json:"message"`
//}
//
////Add new Item the store
//func NewItem(pgx *pgxpool.Pool, ctx context.Context, new_item interface{}) {
//	//aa,ok:=new_item.(Message)
//	info := Message{}
//	mes := new_item.(map[string]interface{})
//	for _, value := range mes {
//		info.message = value.(string)
//	}
//	fmt.Println("new message ", info)
//	ProfileInsert := `INSERT INTO temp(message) VALUES($1);`
//	_, err := pgx.Exec(ctx, ProfileInsert, info.message)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	fmt.Println("data inserted successfully")
//
//}
//
////Update the existing  tem
//func UpdateItem(pgx *pgxpool.Pool, ctx context.Context, item interface{}) {
//	ProfileUpdate := `UPDATE temp SET message =$1;`
//	info := Message{}
//	mes := item.(map[string]interface{})
//	for _, value := range mes {
//		info.message = value.(string)
//	}
//	_, err := pgx.Exec(ctx, ProfileUpdate, info.message)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	fmt.Println("data updated successfully")
//
//}
//
////List all Items from the store
//func GetItems(pgx *pgxpool.Pool, ctx context.Context) []Message {
//	ProfileQuery := `SELECT * FROM temp;`
//	rows, err := pgx.Query(ctx, ProfileQuery)
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil
//	}
//	var information []Message
//	for rows.Next() {
//		info := Message{}
//		err = rows.Scan(&info.message)
//		if err != nil {
//			fmt.Println(err.Error())
//			return nil
//		}
//		information = append(information, info)
//	}
//	defer rows.Close()
//	return information
//}
