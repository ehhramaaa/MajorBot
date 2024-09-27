package core

import (
	"net"
	"net/http"
)

type Client struct {
	account     Account
	proxy       string
	accessToken string
	httpClient  *http.Client
}

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type Account struct {
	queryId        string
	userId         int
	username       string
	firstName      string
	lastName       string
	authDate       string
	hash           string
	allowWriteToPm bool
	languageCode   string
	queryData      string
	walletAddress  string
}
