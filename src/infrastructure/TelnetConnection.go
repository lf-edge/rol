package infrastructure

import (
	"github.com/reiver/go-telnet"
	"rol/app/errors"
	"strings"
)

//TelnetConnection structure for telnet connection
type TelnetConnection struct {
	bond *telnet.Conn
}

//NewTelnetConnection constructor for TelnetConnection
func NewTelnetConnection() *TelnetConnection {
	return &TelnetConnection{}
}

//Connect makes a connection with telnet server
//
//Params:
//	address - telnet server address
//Return:
//	error - if an error occurs, otherwise nil
func (t *TelnetConnection) Connect(address string) (err error) {
	t.bond, err = telnet.DialTo(address)
	if nil != err {
		if err != nil {
			return errors.Internal.Wrap(err, "error connecting to telnet server")
		}
	}
	return nil
}

//Read reads all output lines before expect word
//
//Params:
//	expect - the word to which you want to read lines
//Return:
//	string - telnet output
//	error - if an error occurs, otherwise nil
func (t TelnetConnection) Read(expect string) (string, error) {
	var buffer [1]byte
	recvData := buffer[:]
	var (
		n   int
		err error
		out string
	)
	for {
		n, err = t.bond.Read(recvData)
		if n <= 0 || err != nil || strings.Contains(out, expect) {
			if err != nil {
				return "", errors.Internal.Wrap(err, "error reading from telnet server")
			}
			break
		}
		out += string(recvData)
	}
	return out, nil
}

//Send sends command to telnet server
//
//Params:
//	command - command to send
//Return:
//	error - if an error occurs, otherwise nil
func (t TelnetConnection) Send(command string) error {
	var commandBuffer []byte
	for _, char := range command {
		commandBuffer = append(commandBuffer, byte(char))
	}

	crlfBuffer := [2]byte{'\r', '\n'}
	crlf := crlfBuffer[:]

	_, err := t.bond.Write(commandBuffer)
	if err != nil {
		return errors.Internal.Wrap(err, "error sending command to telnet server")
	}
	_, err = t.bond.Write(crlf)
	if err != nil {
		return errors.Internal.Wrap(err, "error sending line break to telnet server")
	}
	return nil
}
