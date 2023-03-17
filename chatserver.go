package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn     net.Conn
	username string
}

var clients []Client

func main() {
	fmt.Println("Chat Server started...")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("New client connected:", conn.RemoteAddr().String())
	defer conn.Close()
	client := Client{conn: conn}

	client.conn.Write([]byte("Enter your username: "))
	username, err := bufio.NewReader(client.conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	client.username = strings.TrimSpace(username)
	clients = append(clients, client)

	broadcastMessage(client, fmt.Sprintf("%v joined the chat", client.username))
	for {
		msg, err := bufio.NewReader(client.conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			removeClient(client)
			broadcastMessage(client, fmt.Sprintf("%v left the chat", client.username))
			return
		}
		msg = strings.TrimSpace(msg)
		if msg == "/quit" {
			removeClient(client)
			broadcastMessage(client, fmt.Sprintf("%v left the chat", client.username))
			return
		}
		broadcastMessage(client, fmt.Sprintf("%v: %v", client.username, msg))
	}
}

func broadcastMessage(sender Client, message string) {
	for _, client := range clients {
		if sender.conn != client.conn {
			_, err := client.conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func removeClient(client Client) {
	for i, c := range clients {
		if c.conn == client.conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}
