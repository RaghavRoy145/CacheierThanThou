package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

// type Command string

// const (
// 	CMDSet Command = "SET"
// 	CMDGet Command = "GET"
// )

// type Message struct {
// 	Cmd   Command
// 	Key   []byte
// 	Value []byte
// 	TTL   time.Duration
// }

// func (m *Message) ToBytes() []byte {
// 	switch m.Cmd {
// 	case CMDSet:
// 		cmd := fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.TTL)
// 		return []byte(cmd)
// 	case CMDGet:
// 		cmd := fmt.Sprintf("%s %s", m.Cmd, m.Key)
// 		return []byte(cmd)
// 	default:
// 		panic("unknown command")
// 	}
// }

// func parseMessage(rawCmd []byte) (*Message, error) {
// 	var (
// 		rawStr = string(rawCmd)
// 		parts  = strings.Split(rawStr, " ")
// 	)

// 	if len(parts) < 0 {
// 		log.Println("invalid command")
// 		return nil, errors.New("invalid protocol format")
// 	}

// 	msg := &Message{
// 		Cmd: Command(parts[0]),
// 		Key: []byte(parts[1]),
// 	}
// 	if msg.Cmd == CMDSet {
// 		if len(parts) < 4 {
// 			return nil, errors.New("invalid SET format")
// 		}
// 		msg.Value = []byte(parts[2])

// 		ttl, err := strconv.Atoi(parts[3])
// 		if err != nil {
// 			return nil, errors.New("invalid SET TTL")
// 		}
// 		msg.TTL = time.Duration(ttl)
// 	}
// 	return msg, nil
// }

type Command byte

const (
	CmdNonce Command = iota
	CmdSet
	CmdGet
	CmdDel
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

type CommandGet struct {
	Key []byte
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	//read the docs for this again
	binary.Write(buf, binary.LittleEndian, CmdSet)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)

	binary.Write(buf, binary.LittleEndian, int32(c.TTL))

	return buf.Bytes()
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	//read the docs for this again
	binary.Write(buf, binary.LittleEndian, CmdGet)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

func ParseCommand(r io.Reader) any {
	var cmd Command
	binary.Read(r, binary.LittleEndian, &cmd)

	switch cmd {
	case CmdSet:
		return parseSetCommand(r)
	case CmdGet:
		return parseGetCommand(r)
	default:
		panic("invalid command")
	}
}

func parseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Key = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int32
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = int(ttl)
	return cmd
}

func parseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)
	return cmd
}
