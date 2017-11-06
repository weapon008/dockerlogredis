package redislog

import (
	"github.com/docker/docker/daemon/logger"
)

func (rl *Redislog) ReadLogs(cfg logger.ReadConfig) *logger.LogWatcher {
	logWatcher := logger.NewLogWatcher()
	if cfg.Follow {
		rl.bus.Sub(GetGUID(), logWatcher)
	} else {
		go func() {
			msgs := rl.rb.Tail(cfg.Tail, cfg.Since)
			for _, msg := range msgs {
				select {
				case <-logWatcher.WatchClose():
					return
				case logWatcher.Msg <- msg:
				}
			}
		}()
	}
	return logWatcher
}
