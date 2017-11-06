package redislog

import (
	"github.com/docker/docker/daemon/logger"
	"h3d.com/weipeng/dockerlogredis/ringbuffer"
)

const Name = "redis"

type Redislog struct {
	rb  *ringbuffer.RingBuffer
	bus *msgBus
}

func (rl *Redislog) Close() (err error) {
	rl.rb.Close()
	return nil
}
func (rl *Redislog) Name() (name string) {
	return Name
}
func (rl *Redislog) Log(msg *logger.Message) (err error) {
	rl.bus.Pub(msg)
	_, err = rl.rb.Write(msg)
	// TODO send to redis
	return err
}
func New(info logger.Info) (rl logger.Logger, err error) {
	rl = &Redislog{rb: ringbuffer.New(10000), bus: newMsgBus()}
	return rl, nil
}
