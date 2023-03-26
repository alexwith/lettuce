package protocol

import (
	"fmt"
)

func (protocol *RESPProtocol) WriteSimpleString(value string) {
	simpleString := fmt.Sprintf("%s%s\r\n", string(SimpleStringType), value)
	protocol.Connection.Write([]byte(simpleString))
}
