package main

//go:generate ./embed.sh

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

func main() {
	port := flag.String("port", ":9080", "HTTP service port")
	echo := flag.Bool("echo", false, "echo stdin to stdout")

	flag.Parse()

	m := sync.Mutex{}

	lock := func() { m.Lock() }
	unlock := func() { m.Unlock() }

	handlers := make(map[uint64]chan string)
	distribute := func(line string) {
		lock()
		cs := make([]chan string, 0, 16)
		for _, out := range handlers {
			cs = append(cs, out)
		}
		unlock()
		for _, out := range cs {
			// ToDo: Detect and punish slow listeners.
			out <- line
		}
	}

	go func() {
		log.Printf("listening on stdin")
		in := bufio.NewReader(os.Stdin)
		for {
			line, err := in.ReadString('\n')
			if err == io.EOF {
				break
			}
			if *echo {
				fmt.Print(line)
			}
			line = strings.TrimSpace(line)
			distribute(line)
		}
	}()

	i := uint64(0)
	inc := func() uint64 {
		return atomic.AddUint64(&i, 1)
	}

	handler := func(ws *websocket.Conn) {
		id := inc()
		log.Printf("WS client %d", id)
		c := make(chan string, 128)
		lock()
		handlers[id] = c
		unlock()

		go func() {
			for {
				line := <-c
				// ToDo: Maybe detect and reject
				// malformed lines here?  See 'mon.js'
				// for the parser.
				ws.Write([]byte(fmt.Sprintf(`{"line":"%s"}`+"\n", line)))
			}
		}()

		// We don't expect any input, but we'll take a look anyway.
		MaxLine := 128
		in := bufio.NewReaderSize(ws, MaxLine)
		for {
			buf, isPrefix, err := in.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("WS input error %v", err)
				break
			}
			if isPrefix {
				log.Printf("WS input line to long (> %d)", MaxLine)
				break
			}
			log.Printf("WS heard %s", buf)
		}
		log.Printf("leaving WS %d", id)
		if err := ws.Close(); err != nil {
			log.Printf("error %v closing WS %d", err, id)
		}
		lock()
		delete(handlers, id)
		unlock()
	}

	http.Handle("/ws", websocket.Handler(handler))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(monhtml)
	}))

	// Handy for dev/experimentation.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(*port, nil)
	if err != nil {
		panic(err)
	}
}
