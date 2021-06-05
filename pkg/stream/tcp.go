package stream

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/mt-inside/dagger/pkg/value"
)

func Tcp(log logr.Logger, node updateTarget, addr string) {
	log.V(1).Info("stream tcp", "addr", addr)

	// Trivial example; should be a state machine
	// Should also time out to Pending if we don't get a message without a certain period
	for {
		var c net.Conn
		var err error
		for {
			c, err = net.DialTimeout("tcp", addr, 1*time.Second)
			if err == nil {
				break
			}
			log.Error(err, "Can't tcp connect", "addr", addr)
			time.Sleep(1 * time.Second)
		}
		for {
			msg, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				log.Error(err, "Can't read message", "addr", addr)
				break
			}
			n, err := strconv.ParseInt(strings.TrimSpace(string(msg)), 10, 64)
			if err != nil {
				log.Error(err, "Can't parse message", "addr", addr, "message", msg)
				break
			}

			node.Update(value.NewAvailable(n))
			time.Sleep(1 * time.Second)
		}
		c.Close()
		node.Update(value.NewPending())
	}
}
