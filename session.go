// session.go
package main

import (
	"context"
	"math/rand"
	"sync"
)

type TacacsConfig struct {
	IPtype           string //"ip4","ip6"
	ServerIP         string
	ServerPort       uint16
	LocalIP          string
	LocalPort        uint16
	ConnMultiplexing bool
}

func ConfigSet(config TacacsCOnfig) {
	TacacsMng.Lock()
	defer TacacsMng.Unlock()
	TacacsMng.Config = config
}

func ConfigGet() (config TacacsCOnfig) {
	TacacsMng.Lock()
	defer TacacsMng.Unlock()
	config := TacacsMng.Config
	return config
}

const (
	MaxUint8 = ^uint8(0)
)

type Config struct {
}

type Manager struct {
	Sessions sync.Map

	Trans *Transport
	ctx   context.Context
	sync.RWMutex
	Config TacacsConfig
}

type Session struct {
	sync.Mutex
	SessionSeqNo uint8
	SessionID    uint32
	UserName     string
	Password     string
	ReadBuffer   chan []byte
	mng          *Manager
	t            *Transport
	ctx          context.Context
}

func NewSession(ctx context.Context, name, passwd string) *Session {
	sess := &Session{}
	sess.Password = passwd
	sess.UserName = name
	sess.SessionSeqNo = 1
	sess.ReadBuffer = make(chan []byte, 10)
	sess.mng = TacacsMng
	sess.ctx = ctx
	rand.Seed(time.Now().Unix())
	SessionID := rand.Uint32()
	for {
		if _, ok := TacacsMng.Sessions.Load(SessionID); ok {
			SessionID = rand.Uint32()
		} else {
			break
		}
	}
	sess.SessionID = SessionID
	sess.t = newTransport(ctx, TacacsMng.Config)

	TacacsMng.Sessions.Store(SessionID, &sess)
	return sess
}

var TacacsMng *Manager

func TacacsInit() {
	TacacsMng = &Manager{}
	TacacsMng.ctx = context.TODO()
}

func (sess *Session) close() {
	sess.mng.Sessions.Delete(sess.SessionID)
	sess.t.close()
	//close(sess.)
}
