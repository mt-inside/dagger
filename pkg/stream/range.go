package stream

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/mt-inside/dagger/pkg/value"
)

func Range(log logr.Logger, node updateTarget, step int64, period time.Duration) {
	log.V(1).Info("stream range", "step", step, "period", period)

	n := int64(0)
	for {
		node.Update(value.NewAvailable(n))
		n += step
		time.Sleep(period)
	}
}
