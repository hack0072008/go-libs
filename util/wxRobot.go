package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

/*
 useage:
   1.file = NewWeiXinRobot(path)
   2. file.open()
   3.RegisterDefaultFile(type, file)
 */

type message struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"markdown"`
}
type WeiXinRobot struct {
	url string
	sync.Mutex
}

func NewWeiXinRobot(url string) (BackEndPush, error) {
	return &WeiXinRobot{
		url:   url,
		Mutex: sync.Mutex{},
	}, nil
}

func (w *WeiXinRobot) Open() error {
	_, err := url.Parse(w.url)
	return err
}

func (w *WeiXinRobot) Push(b []byte) error {
	w.Lock()
	defer w.Unlock()
	var pushMessage = message{
		MsgType: "markdown",
		Text: struct {
			Content string `json:"content"`
		}{string(b)},
	}
	messageBody, _ := json.Marshal(pushMessage)
	req, err := http.NewRequest("POST", w.url, bytes.NewReader(messageBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}

func (w *WeiXinRobot) Close() error {
	return nil
}


// weixin interface
type BackEndPush interface {
	Open() error
	Push(b []byte) error
	Close() error
}

const (
	Default = "default"
	WXRobot = "WeiXinRobot"
)

var (
	RegistryBackEnd map[string]BackEndPush
	RegistryLock    sync.RWMutex
)

func init() {
	RegistryBackEnd = map[string]BackEndPush{}
}

func RegisterDefaultFile(name string, push BackEndPush) error {
	RegistryLock.Lock()
	defer RegistryLock.Unlock()
	RegistryBackEnd[name] = push
	return nil
}

func GetBackEnd(backend string) (BackEndPush, error) {
	RegistryLock.RLock()
	defer RegistryLock.RUnlock()
	push, ok := RegistryBackEnd[backend]
	if !ok {
		return nil, fmt.Errorf("Not Find BackEend Push Server")
	}
	return push, nil
}

