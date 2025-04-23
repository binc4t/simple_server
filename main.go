package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func HandleMockFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	content, err := os.ReadFile("test.ts")
	if err != nil {
		fmt.Println("read file err, ", err)
	}

	// set linger
	conn := GetTCPConnFromContext(r.Context())
	err = conn.SetLinger(5)
	if err != nil {
		fmt.Println("set linger err, ", err)
	} else {
		fmt.Println("set linger success")
	}

	// write content length
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))

	_, err = w.Write(content)
	if err != nil {
		fmt.Printf("time: %s, http write err, %s\n", time.Now().Format("2006-01-02 15:04:05"), err)
	}
	fmt.Printf("time: %s, done, time_cost: %s\n", time.Now().Format("2006-01-02 15:04:05"), time.Since(start))
}

func GetConnFromContext(ctx context.Context) net.Conn {
	return ctx.Value("conn").(net.Conn)
}

func GetTCPConnFromContext(ctx context.Context) *TCPConn {
	return ctx.Value("conn").(*TCPConn)
}

func main() {
	server := http.Server{
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, "conn", c)
		},
		Addr: ":8080",
	}

	http.Handle("/mock", http.HandlerFunc(HandleMockFile))
	fmt.Println("Listen on port 8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
