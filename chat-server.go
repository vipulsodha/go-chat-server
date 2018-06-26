package main

import (
	"golang.org/x/net/websocket"
	"net/http"
	"log"
	"sync"
	"fmt"
	"bufio"
	"os"
)

var socketMap = make(map[int]*websocket.Conn)

var socketCount int = 0;

func main()  {

	http.Handle("/", websocket.Handler(handleWs))

	go readMessageConsole()

	log.Fatal(http.ListenAndServe(":8088", nil))

}

/**
For sending message from terminal directly
 */
func readMessageConsole()  {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		broadCastMessageToAll(scanner.Text(), -1)
	}

}

/**
Handling WS
 */
func handleWs(ws *websocket.Conn)  {

	var m sync.Mutex;

	var wait sync.WaitGroup;

	wait.Add(1)


	ch := make(chan int)

	go addSocketToMap(ws, &m, ch)

	key := <- ch


	go listenToWs(ws, key, &wait)

	//fmt.Println("New socket added success");

	wait.Wait()

}

/**
Maintaining map for Websockets that are online
 */
func addSocketToMap(ws *websocket.Conn, m *sync.Mutex, ch chan int)  {

	m.Lock()

	socketMap[socketCount] =  ws;
	ch <- socketCount
	socketCount++

	m.Unlock()

}

/**
Listening to websockets for new messages
 */
func listenToWs(ws *websocket.Conn, key int, wait *sync.WaitGroup)  {

	for  {

		var reply string;

		websocket.Message.Receive(ws, &reply)

		fmt.Println(reply)

		go broadCastMessageToAll(reply, key)

	}

	wait.Done()

}

/**
Broadcast message to all the online sockets
 */
func broadCastMessageToAll(message string, socketKey int)  {

	for key, ws := range socketMap {

		if key == socketKey {
			continue
		}
		websocket.Message.Send(ws, message)
	}
}