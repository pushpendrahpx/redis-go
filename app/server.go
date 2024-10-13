package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit


func do(c net.Conn) {

	defer closeConnection(c)


	buf := make([]byte, 128)

	n, err := c.Read(buf)
	fmt.Printf("Read %d chars\n", n)
	if err != nil {
		fmt.Errorf("Got Error : %v", err)
	}

	fmt.Println(buf, string(buf))

	n, err = c.Write([]byte("+PONG\r\n"))

	fmt.Printf("Wrote %d chars\n", n)
	fmt.Println(c.RemoteAddr())
	if err != nil {
		fmt.Errorf("Got Error while writing : %v", err)
	}
	


}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {

		
		// This is a Blcking Call
		c, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}else{
			go do(c)
		}
	}
}


func closeConnection(c net.Conn) {
	err := c.Close()

	if err != nil {
		fmt.Errorf("Error: %v", err)
		panic("Connection Closed")
	}
}