package protocol

func (protocol *RESPProtocol) WriteSimpleString(value string) {
	protocol.WriteRawString(protocol.CreateSimpleString(value))
}

func (protocol *RESPProtocol) WriteError(value string) {
	protocol.WriteRawString(protocol.CreateError(value))
}

func (protocol *RESPProtocol) WriteInteger(value int) {
	protocol.WriteRawString(protocol.CreateInteger(value))
}

func (protocol *RESPProtocol) WriteBulkString(value string) {
	protocol.WriteRawString(protocol.CreateBulkString(value))
}

func (protocol *RESPProtocol) WriteArray(value interface{}) {
	protocol.WriteRawString(protocol.CreateArray(value))
}

func (protocol *RESPProtocol) WriteRawString(value string) {
	protocol.Connection.Write([]byte(value))
}
