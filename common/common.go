package common

import (
	"net"
	"sync"
)

// UserID
type UserID int64

// Server Config
var (
	// Development Mode.
	AppMode string

	// Server Listener Addr.
	RequestAddr  string
	ResponseAddr string
	ChatAddr     string
	FileAddr     string
	VideoAddr    string

	// DataBase Config.
	DbHost string
	DbUser string
	DbPort string
	DbName string
	DbPass string
)

// All Listeners of an user.
var (
	RequestListener  net.Listener
	ResponseListener net.Listener
	ChatListener     net.Listener
	FileListener     net.Listener
	VideoListener    net.Listener
)

// All Connections of an user.
type Conns struct {
	RequestConn  net.Conn
	ResponseConn net.Conn
	ChatConn     net.Conn
	FileConn     net.Conn
	VideoConn    net.Conn
}

// Aes Key to encrypt password.
var AESKEY string

// To save All Login Users.
var AllUsers sync.Map
