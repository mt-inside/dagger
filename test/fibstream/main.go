package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	var val int64

	go func() {
		a := int64(1)
		b := int64(1)
		for {
			val = a + b
			b = a + b
			a = b - a
			if a+b > 1000000 {
				a = 1
				b = 1
			}
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Returning %d", val)
		fmt.Fprintf(w, strconv.FormatInt(val, 10))
	})
	log.Fatal(http.ListenAndServe(os.Args[1], nil))
}
