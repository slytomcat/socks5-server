package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/slytomcat/go-socks5"
)

const (
	UserEnv     = "PROXY_USER"
	PasswordEnv = "PROXY_PASSWORD"
	PortEnv     = "PROXY_PORT"
)

func main() {

	port := os.Getenv(PortEnv)
	if port == "" {
		port = "1080"
	}

	//Initialize socks5 config
	socks5conf := &socks5.Config{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	if user, password := os.Getenv(UserEnv), os.Getenv(PasswordEnv); len(user)+len(password) != 0 {
		creds := socks5.StaticCredentials{user: password}
		socks5conf.AuthMethods = []socks5.Authenticator{socks5.UserPassAuthenticator{Credentials: creds}}
	}

	server, err := socks5.New(socks5conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Start listening proxy service on port %s\n", port)
	if err := server.ListenAndServe("tcp", ":"+port); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
}
