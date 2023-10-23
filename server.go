package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/raghavroy145/DistributedCaching/cache"
	"github.com/raghavroy145/DistributedCaching/proto"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}
	log.Printf("server on starting on port [%s]\n", s.ListenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	// buf := make([]byte, 2048)
	// fmt.Println("connection made:", conn.RemoteAddr())
	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error:", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}
	// fmt.Println("connection closed: ", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {

	switch v := cmd.(type) {
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:
		s.handleGetCommand(conn, v)
	}
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet) error {
	// log.Printf("GET %s", cmd.Key)

	resp := proto.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}
	resp.Status = proto.StatusOk
	resp.Value = value
	_, err = conn.Write(resp.Bytes())
	return err
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	// log.Printf("SET %s to %s", cmd.Key, cmd.Value)
	resp := proto.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}
	resp.Status = proto.StatusOk
	_, err := conn.Write(resp.Bytes())
	return err
}

// func respondClient(conn net.Conn, msg any) error {
// 	// _, err := conn.Write(msg.Bytes())
// }
