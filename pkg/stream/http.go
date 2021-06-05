package stream

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/mt-inside/dagger/pkg/value"
)

func Http(log logr.Logger, node updateTarget, url string) {
	log.V(1).Info("stream url", "url", url)

	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Error(err, "Can't fetch stream", "url", url)
			os.Exit(1)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "Can't read body", "url", url)
			os.Exit(1)
		}
		n, err := strconv.ParseInt(string(body), 10, 64)
		if err != nil {
			log.Error(err, "Can't parse body", "url", url, "body", body)
			os.Exit(1)
		}

		node.Update(value.NewAvailable(n))
		time.Sleep(1 * time.Second)
	}
}
