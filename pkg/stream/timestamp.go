package stream

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/mt-inside/dagger/pkg/value"
)

func Timestamp(log logr.Logger, node updateTarget, period time.Duration) {
	log.V(1).Info("stream timestamp", "period", period)

	for {
		node.Update(value.NewAvailable(time.Now(time.Now().Unix()))
		time.Sleep(period)
	}
}
