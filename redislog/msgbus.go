package redislog

import (
	"github.com/docker/docker/daemon/logger"
)

type msgBus struct {
	watchers map[string]*logger.LogWatcher
}

func (bus *msgBus) Sub(key string, lw *logger.LogWatcher) {
	bus.watchers[key] = lw
}
func (bus *msgBus) Pub(msg *logger.Message) {
	for key, lw := range bus.watchers {
		select {
		case <-lw.WatchClose():
			delete(bus.watchers, key)
		case lw.Msg <- msg:
		}
	}
}
func newMsgBus() (bus *msgBus) {
	bus = &msgBus{watchers: make(map[string]*logger.LogWatcher)}
	return bus
}
