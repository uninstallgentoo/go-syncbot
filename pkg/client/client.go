package client

import (
	"fmt"
	"log"
	"time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"go.uber.org/zap"

	"sync-bot/pkg/config"
	"sync-bot/pkg/models"
)

//socket.io events
const (
	JoinChannelMethod = "joinChannel"
	LoginMethod       = "login"
	ChatMessage       = "chatMsg"
	OnSetUserRank     = "setUserRank"
	OnAddUser         = "addUser"
	OnUserlist        = "userlist"
	OnUserLeave       = "userLeave"
	OnConnection      = gosocketio.OnConnection
	OnDisconnection   = gosocketio.OnDisconnection
)

var now = time.Now().Unix() * 1000

type Channel struct {
	Name     string `json:"name"`
	Password string `json:"pw"`
}

type SocketClient struct {
	conf   *config.Config
	conn   *gosocketio.Client
	logger *zap.SugaredLogger
	chat   Chat
}

func NewSocketClient(conf *config.Config, chat *Chat, logger *zap.SugaredLogger) *SocketClient {
	return &SocketClient{
		conf:   conf,
		chat:   *chat,
		logger: logger,
	}
}

func (s *SocketClient) connect() {
	client, err := gosocketio.Dial(
		gosocketio.GetUrl(s.conf.Host, s.conf.Port, s.conf.Secure),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		log.Fatalf("Error has occured during conect to socket.io: %e", err)
	}
	s.conn = client
}

func (s *SocketClient) joinChannel() {
	var err error
	err = s.join()
	err = s.login()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SocketClient) onConnection() {
	err := s.conn.On(OnConnection, func(c *gosocketio.Channel) {
		fmt.Printf("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SocketClient) onDisconnection() {
	err := s.conn.On(OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected")
	})
	if err != nil {
		s.logger.Fatal(err)
	}
}

func (s *SocketClient) onAddUser(_ *gosocketio.Channel, user *models.User) {
	err := s.chat.userService.SaveNewUser(user)
	if err != nil {
		print(err)
	}
	_, ok := s.chat.users[user.Name]
	if !ok {
		s.chat.users[user.Name] = user
	}
}

func (s *SocketClient) onUserLeave(_ *gosocketio.Channel, user *models.UserLeave) {
	_, ok := s.chat.users[user.Name]
	if ok {
		delete(s.chat.users, user.Name)
	}
}

func (s *SocketClient) onSetUserRank(_ *gosocketio.Channel, user *models.UpdatedUser) {
	err := s.chat.userService.UpdateUserRank(user)
	if err != nil {
		print(err)
	}
}

func (s *SocketClient) onUserList(_ *gosocketio.Channel, users []*models.User) {
	for _, user := range users {
		s.chat.users[user.Name] = user
	}
}

func (s *SocketClient) join() error {
	return s.conn.Emit(JoinChannelMethod, Channel{s.conf.Channel.Name, s.conf.Channel.Password})
}

func (s *SocketClient) login() error {
	return s.conn.Emit(LoginMethod, Channel{s.conf.User.Name, s.conf.User.Password})
}

func (s *SocketClient) onChatMessage(_ *gosocketio.Channel, msg *models.Message) {
	msgTimestamp := msg.Time
	if msgTimestamp < now {
		return
	}
	if len(s.chat.messages) == 100 {
		err := s.chat.chatService.SaveChatHistory(s.chat.messages)
		if err != nil {
			s.logger.Error(err)
		} else {
			s.chat.messages = make([]*models.Message, 0, 100)
		}
	}
	go s.chat.handleMessage(msg)
	s.Emit()
}

func (s *SocketClient) registerEvents() {
	var err error
	err = s.conn.On(OnConnection, s.onConnection)
	err = s.conn.On(OnDisconnection, s.onDisconnection)
	err = s.conn.On(ChatMessage, s.onChatMessage)
	err = s.conn.On(OnAddUser, s.onAddUser)
	err = s.conn.On(OnSetUserRank, s.onSetUserRank)
	err = s.conn.On(OnUserlist, s.onUserList)
	err = s.conn.On(OnUserLeave, s.onUserLeave)

	if err != nil {
		s.logger.Fatal(err)
	}
}

func (s *SocketClient) Emit() {
	for m := range s.chat.commandResults {
		err := s.conn.Emit(m.Method, m.Message)
		if err != nil {
			s.logger.Error(err)
		}
	}
}

func (s *SocketClient) Start() {
	s.connect()
	s.logger.Info("Connected to socket.io")
	s.joinChannel()
	s.logger.Info("Join to channel")
	s.registerEvents()
	done := make(chan bool)
	go func() {
		for {
			time.Sleep(30 * time.Second)
			done <- s.conn.IsAlive()
		}
	}()

	for msg := range done {
		if !msg {
			s.logger.Fatal("Server shutdown")
			return
		}
	}
}
