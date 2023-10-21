package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

type Message struct {
	Cmd   Command
	Key   []byte
	Value []byte
	TTL   time.Duration
}

func parseMessage(rawCmd []byte) (*Message, error) {
	var (
		rawStr = string(rawCmd)
		parts  = strings.Split(rawStr, " ")
	)

	if len(parts) < 0 {
		log.Println("invalid command")
		return nil, errors.New("invalid protocol format")
	}

	msg := &Message{
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}
	if msg.Cmd == CMDSet {
		if len(parts) < 4 {
			return nil, errors.New("invalid SET format")
		}
		msg.Value = []byte(parts[2])

		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("invalid SET TTL")
		}
		msg.TTL = time.Duration(ttl)
	}
	return msg, nil
}
