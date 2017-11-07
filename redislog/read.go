package redislog

import (
	"bytes"
	"io"

	"github.com/docker/docker/daemon/logger"
)

func (rl *Redislog) ReadLogs(cfg logger.ReadConfig) *logger.LogWatcher {
	logWatcher := logger.NewLogWatcher()

	if cfg.Follow {
		rl.bus.Sub(GetGUID(), logWatcher)
	} else {
		go func() {

			msgs := rl.rb.Tail(cfg.Tail, cfg.Since)
			i := 0
			for {

				if i >= len(msgs) {
					logWatcher.Err <- io.EOF
					break
				}
				msg := msgs[i]
				msg.Line = append(bytes.TrimSpace(msg.Line), []byte("\n")...)

				i += 1
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
