package connection

import (
	"fmt"
	"net"
	"os"
)

const CONNECTION_TYPE = "tcp"

func Listen(host string, port int16) {
	listener, err := net.Listen(CONNECTION_TYPE, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Failed to start the lettuce server")
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Printf("The lettuce server has started on %s:%d\n", host, port)

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err.Error())
			continue
		}

		HandleConnection(connection)
	}
}
