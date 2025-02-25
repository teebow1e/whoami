package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	port    string
	name    string
	verbose bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.StringVar(&port, "port", getEnv("WHOAMI_PORT_NUMBER", "8080"), "give me a port number")
	flag.StringVar(&name, "name", os.Getenv("WHOAMI_NAME"), "give me a name")
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer listener.Close()

	log.Printf("Starting TCP whoami server on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(30 * time.Second))

	initialResponse := getWhoamiResponse(conn)
	_, err := conn.Write([]byte(initialResponse + "\n> "))
	if err != nil {
		log.Printf("Error writing initial response: %v", err)
		return
	}

	buffer := make([]byte, 1024)
	for {
		conn.SetDeadline(time.Now().Add(30 * time.Second))

		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v", err)
			}
			return
		}
		command := strings.TrimSpace(string(buffer[:n]))

		var response string
		switch {
		case command == "whoami":
			response = getWhoamiResponse(conn)

		case command == "bench":
			response = "1"

		case strings.HasPrefix(command, "/data?size="):
			sizeStr := strings.TrimSpace(strings.TrimPrefix(command, "/data?size="))
			size, err := parseSize(sizeStr)
			if err != nil {
				response = fmt.Sprintf("Error: invalid size parameter: %v", err)
			} else {
				response = generateData(size)
			}

		case command == "":
			response = ""

		default:
			response = "that command doesnt exist, but the server is still functioning good!"
		}

		if verbose && command != "" {
			log.Printf("%s - [%s] command: %s",
				conn.RemoteAddr().String(),
				time.Now().Format(time.RFC1123),
				command)
		}

		if response != "" {
			_, err = conn.Write([]byte(response + "\n> "))
		} else {
			_, err = conn.Write([]byte("> "))
		}
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			return
		}
	}
}

func getWhoamiResponse(conn net.Conn) string {
	var response strings.Builder

	response.WriteString("Yo, skibidi bop bop, what's good, fam!\n")
	response.WriteString("If you're clockin' these lines, your TCP hookup with the remote server is skibidi poppin', no cap, connection's straight fire.\n")
	response.WriteString("Here's the drip on that link-up, bop bop:\n")
	response.WriteString("\n")

	if name != "" {
		response.WriteString(fmt.Sprintf("Name: %s\n", name))
	}

	hostname, err := os.Hostname()
	if err == nil {
		response.WriteString(fmt.Sprintf("Hostname: %s\n", hostname))
	}

	remoteAddr := conn.RemoteAddr().String()
	response.WriteString(fmt.Sprintf("RemoteAddr: %s\n", remoteAddr))
	response.WriteString(fmt.Sprintf("Time: %s\n", time.Now().Format(time.RFC1123)))

	return response.String()
}

func parseSize(sizeStr string) (int64, error) {
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}
	if size < 0 {
		return 0, fmt.Errorf("size must be non-negative")
	}
	return size, nil
}

func generateData(size int64) string {
	if size == 0 {
		return ""
	}

	var builder strings.Builder
	const charset = "-ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Write initial character
	builder.WriteByte('|')

	// Write content
	for i := int64(1); i < size-1; i++ {
		builder.WriteByte(charset[i%int64(len(charset))])
	}

	// Write final character if size > 1
	if size > 1 {
		builder.WriteByte('|')
	}

	return builder.String()
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
