package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Yideg/webSockExample/entity"
	"github.com/Yideg/webSockExample/global"
	"github.com/Yideg/webSockExample/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/pgxpool"
	"gopkg.in/olahol/melody.v1"
	"log"
)

//var clients = make(map[string]interface{})
//var sessions []*melody.Session

type WebSocket struct {
	M    *melody.Melody
	Repo *repository.PgxRepo
}

func New(mm *melody.Melody, r *repository.PgxRepo) WebSocket {
	return WebSocket{M: mm, Repo: r}
}

func (m WebSocket) ConvertToWebsocketConnect(c *gin.Context) {
	err := m.M.HandleRequest(c.Writer, c.Request)
	fmt.Println("error", err)
	if err != nil {
		log.Fatal("Unable to convert the http connection to websocket")
		return
	}
}

func (m WebSocket) HandleDisconnect() {
	m.M.HandleDisconnect(func(s *melody.Session) {
		for index, ss := range global.Sessions {
			if s == ss {
				global.Sessions = append(global.Sessions[:index], global.Sessions[(index+1):]...)
			}
		}
		fmt.Println("no sessions formed before is ", len(global.Sessions))
	})
}

func (m WebSocket) HandleConnect() {
	m.M.HandleConnect(func(s *melody.Session) {
		global.Sessions = append(global.Sessions, s)
		fmt.Println("New session is Added to Session collections", len(global.Sessions))
		fmt.Println("connecting ...")
		host := s.Request.Host + s.Request.URL.Path
		key := s.Request.URL.Path
		s.Set(key, host)
		fmt.Println("number of  session ", len(global.Sessions))

	})
}
func (m WebSocket) SendMessage(p *pgxpool.Pool, c context.Context) {
	m.M.HandleMessage(func(s *melody.Session, msg []byte) {
		messages := m.Repo.GetItems(p, c)
		var mess interface{}
		err := json.Unmarshal(msg, &mess)
		if err != nil {
			log.Fatal("error ", err)
			return
		}

		info := entity.Message{}
		mes := mess.(map[string]interface{})
		for _, value := range mes {
			info.Text = value.(string)
		}
		messages = append(messages, info)
		var ses []*melody.Session
		current_connected := s.Keys[s.Request.URL.Path]
		for _, se := range global.Sessions {
			find := se.Keys[s.Request.URL.Path]
			if find == current_connected {
				ses = append(ses, se)
			}
		}

		for i := 0; i < len(messages); i++ {
			mm := messages[i]
			err := m.M.BroadcastMultiple([]byte(mm.Text), ses)
			if err != nil {
				log.Fatal("error has been occurred during sending messages")
				return
			}
		}
		//insert values into the table
		m.Repo.NewItem(p, c, mes)

	})

}
