package client

import (
	"github.com/gorcon/rcon"
)

type RCON struct {
	conn *rcon.Conn
}

func Dial(address string, password string) (*RCON, error) {
	conn, err := rcon.Dial(address, password)
	if err != nil {
		return nil, err
	}
	return &RCON{conn: conn}, nil
}

func (r *RCON) Close() error {
	return r.conn.Close()
}

func (r *RCON) Execute(command string) (string, error) {
	return r.conn.Execute(command)
}
