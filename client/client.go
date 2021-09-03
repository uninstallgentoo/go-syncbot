package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"go.uber.org/zap"

	"sync-bot/config"
	"sync-bot/models"
	"sync-bot/processors"
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
)

type Server struct {
	URL    string `json:"url"`
	Secure bool   `json:"secure"`
}
type ServerConfig struct {
	Servers []Server `json:"servers"`
}

type Channel struct {
	Name     string `json:"name"`
	Password string `json:"pw"`
}

type Events map[string]interface{}

type SocketClient struct {
	conf      *config.Config
	conn      *gosocketio.Client
	logger    *zap.Logger
	chat      processors.ChatHandler
	commands  processors.CommandHandler
	startedAt int64
}

// TODO: logger, pass ctx to repository
func NewSocketClient(conf *config.Config, chat processors.ChatHandler, cmdHandler processors.CommandHandler, logger *zap.Logger) *SocketClient {
	return &SocketClient{
		conf:      conf,
		chat:      chat,
		commands:  cmdHandler,
		logger:    logger,
		startedAt: time.Now().Unix() * 1000,
	}
}

func (s *SocketClient) initEventsMap() Events {
	return Events{
		ChatMessage:   s.onChatMessage,
		OnAddUser:     s.onAddUser,
		OnSetUserRank: s.onSetUserRank,
		OnUserlist:    s.onUserList,
		OnUserLeave:   s.onUserLeave,
	}
}

func (s *SocketClient) connect() (*gosocketio.Client, error) {
	sioConfig, err := s.fetchServerConfig()
	if err != nil {
		return nil, err
	}
	serverConfig := &ServerConfig{}
	if err := json.Unmarshal(sioConfig, serverConfig); err != nil {
		return nil, err
	}
	var host string
	var port int
	var secure bool
	for _, item := range serverConfig.Servers {
		if s.conf.Secure == item.Secure {
			// URL example: "http://somehost.cytu.be:3000",
			portStartIdx := strings.LastIndex(item.URL, ":")
			hostStartIdx := strings.LastIndex(item.URL, "/")
			host = item.URL[hostStartIdx+1 : portStartIdx]
			secure = item.Secure
			port, err = strconv.Atoi(item.URL[portStartIdx+1:])
			if err != nil {
				return nil, err
			}
		}
	}
	client, err := gosocketio.Dial(
		gosocketio.GetUrl(host, port, secure),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		return nil, err
	}
	return client, err
}

func (s *SocketClient) fetchServerConfig() ([]byte, error) {
	protocol := "http"
	if s.conf.Secure {
		protocol = "https"
	}
	resp, err := http.Get(fmt.Sprintf("%s://%s/socketconfig/%s.json", protocol, s.conf.Host, s.conf.Channel.Name))
	if err != nil {
		return nil, fmt.Errorf("server respond with error during request config: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error has occured during reading server config response : %v", err)
	}
	return body, err
}

func (s *SocketClient) joinChannel() error {
	if err := s.join(); err != nil {
		return err
	}
	if err := s.login(); err != nil {
		return err
	}
	return nil
}

func (s *SocketClient) onAddUser(_ *gosocketio.Channel, user models.User) {
	if err := s.chat.SaveNewUser(user); err != nil {
		s.logger.Error("fail to save new user", zap.Error(err))
	}
}

func (s *SocketClient) onUserLeave(_ *gosocketio.Channel, user models.UserLeave) {
	s.chat.DeleteUser(user)
}

func (s *SocketClient) onSetUserRank(_ *gosocketio.Channel, user models.UpdatedUser) error {
	return s.chat.UpdateUserRank(user)
}

func (s *SocketClient) onUserList(_ *gosocketio.Channel, users []models.User) {
	s.chat.AddUserToList(users)
}

func (s *SocketClient) join() error {
	return s.conn.Emit(JoinChannelMethod, Channel{s.conf.Channel.Name, s.conf.Channel.Password})
}

func (s *SocketClient) login() error {
	return s.conn.Emit(LoginMethod, Channel{s.conf.User.Name, s.conf.User.Password})
}

func (s *SocketClient) onChatMessage(_ *gosocketio.Channel, msg models.Message) {
	msgTimestamp := msg.Time
	if msgTimestamp < s.startedAt {
		return
	}

	s.chat.HandleMessage(msg)
	s.commands.Handle(msg)
	if err := s.Emit(); err != nil {
		s.logger.Error("fail send message to socket.io", zap.Error(err))
	}
}

func (s *SocketClient) registerEvents() error {
	events := s.initEventsMap()
	for k, v := range events {
		err := s.conn.On(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SocketClient) Emit() error {
	for m := range s.chat.GetCommandResults() {
		if err := s.conn.Emit(m.Method, m.Message); err != nil {
			return err
		}
	}
	return nil
}

func (s *SocketClient) Start() {
	connection, err := s.connect()
	if err != nil {
		s.logger.Fatal("fail to connect to socker.io server", zap.Error(err))
	}
	s.conn = connection
	s.logger.Info("Connected to socket.io")
	if err := s.joinChannel(); err != nil {
		s.logger.Fatal("fail to join channel", zap.Error(err))
	}
	s.logger.Info("Join to channel")
	if err := s.registerEvents(); err != nil {
		s.logger.Fatal("fail to subscription to event", zap.Error(err))
	}
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
