package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type GlobalMapValue struct {
	value string
	exp time.Time
}
var globalMap = map[string]GlobalMapValue{}

var CONFIG = map[string]string{}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	for i, flag := range os.Args {
		if flag == "--dir" && i+1 < len(os.Args) {
			CONFIG["dir"] = os.Args[i+1]
		}else if flag == "--dbfilename" && i+1 < len(os.Args)  {
			CONFIG["dbfilename"] = os.Args[i+1]
		}
	}
	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		// Listen for connections
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go runConnection(conn)
		
	}
}

func runConnection(conn net.Conn) {
	
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		command := value.array[0].bulk
		fmt.Println(value.array, len(value.array))
		if command == "ECHO" {
			conn.Write([]byte("+"+value.array[1].bulk+"\r\n"))
		}else if command == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		}else if command == "SET" {
			key := value.array[1].bulk
			v := value.array[2].bulk
			expiry := time.Now()
			if len(value.array) >= 5 && value.array[3].bulk == "px" {
				// means there may be PX <time>
				delay, err := strconv.Atoi(value.array[4].bulk)
				if err != nil {
					conn.Write([]byte("+"+"FailedPX"+"\r\n"))
					return
				}

				expiry = time.Now().Add(time.Duration(delay)*time.Millisecond)
			}else{
				expiry = time.Time{}
			}
			
			globalMap[key] = GlobalMapValue{value: v, exp: expiry}
			conn.Write([]byte("+OK\r\n"))
		}else if command == "GET" {
			key := value.array[1].bulk
			v, err := globalMap[key]
			if err == false {
				conn.Write([]byte("$-1\r\n"))
				continue
			}
			current := time.Now()
			if v.exp.IsZero() {
				conn.Write([]byte("+"+v.value+"\r\n"))
				continue
			}
			if current.UnixMicro() > v.exp.UnixMicro() {
				conn.Write([]byte("$-1\r\n"))
				continue
			}
			conn.Write([]byte("+"+v.value+"\r\n"))
			
			
			
			
		}else if command == "CONFIG" {
			key := value.array[2].bulk
			if key == "dir" {
				dir := CONFIG["dir"]
				conn.Write([]byte("*2\r\n$3\r\ndir\r\n$"+strconv.Itoa(len(dir))+"\r\n"+dir+"\r\n"))
				
			}else if key == "dbfilename" {
				dbfilename := CONFIG["dbfilename"]
				conn.Write([]byte("*2\r\n$3\r\ndbfilename\r\n$"+strconv.Itoa(len(dbfilename))+"\r\n"+dbfilename+"\r\n"))
			}else{
				conn.Write([]byte("+-1\r\n"))
			}
		}else{
			conn.Write([]byte("+OK\r\n"))
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