package redislog

import (
	"log"

	"github.com/docker/docker/daemon/logger"
)

func (rl *Redislog) ReadLogs(cfg logger.ReadConfig) *logger.LogWatcher {
	logWatcher := logger.NewLogWatcher()
	log.Println(`!!!!!!!!!!111`, cfg.Follow)
	if cfg.Follow {
		rl.bus.Sub(GetGUID(), logWatcher)
	} else {
		go func() {
			log.Println(`!!!!!!!!!!222`)
			msgs := rl.rb.Tail(cfg.Tail, cfg.Since)
			i := 0
			for {
				log.Println(`!!!!!!!!!!333`, i)

				if i >= len(msgs) {

					break
				}
				msg := msgs[i]
				i += 1
				select {
				case <-logWatcher.WatchClose():
					return
				case logWatcher.Msg <- msg:
				}
			}
			log.Println(`!!!!!!!!!!888`)

		}()
	}
	return logWatcher
}
