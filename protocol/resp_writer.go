package protocol

func (connection *Connection) WriteSimpleString(value string) {
	connection.WriteRawString(connection.CreateSimpleString(value))
}

func (connection *Connection) WriteError(value string) {
	connection.WriteRawString(connection.CreateError(value))
}

func (connection *Connection) WriteInteger(value int) {
	connection.WriteRawString(connection.CreateInteger(value))
}

func (connection *Connection) WriteBulkString(value string) {
	connection.WriteRawString(connection.CreateBulkString(value))
}

func (connection *Connection) WriteArray(value []any) {
	connection.WriteRawString(connection.CreateArray(value))
}

func (connection *Connection) WriteNullBulkString() {
	connection.WriteRawString(NullBulkString)
}

func (connection *Connection) WriteNullArray() {
	connection.WriteRawString(NullArray)
}

func (connection *Connection) WriteRawString(value string) {
	connection.Handle.Write([]byte(value))
}
