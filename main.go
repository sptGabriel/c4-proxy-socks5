package main

import (
	"github.com/sptGabriel/socks5/app"
)

func main() {
	// c, err := app.NewClient("https://www.proxyscan.io/api", 15*time.Second)
	// if err != nil {
	// 	log.Fatal("cannot create client")
	// }

	srv := app.New()
	// srv.AuthNoAuthenticationRequiredCallback = func(conn *app.Conn) error { return nil }
	srv.ListenAndServe("127.0.0.1:1080")
}
