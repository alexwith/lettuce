package protocol

import (
	"fmt"
)

func (protocol *RESPProtocol) WriteSimpleString(value string) {
	simpleString := fmt.Sprintf("%s%s\r\n", string(SimpleStringType), value)
	protocol.Connection.Write([]byte(simpleString))
}

func (protocol *RESPProtocol) WriteError(value string) {
	error := fmt.Sprintf("%s%s\r\n", string(ErrorType), value)
	protocol.Connection.Write([]byte(error))
}