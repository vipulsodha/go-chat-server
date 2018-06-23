package main

import (
	"golang.org/x/net/websocket"
	"net/http"
	"log"
	"sync"
	"fmt"
)

var socketMap = make(map[int]*websocket.Conn)

var socketCount int = 0;

func main()  {

	http.Handle("/", websocket.Handler(handleWs))

	log.Fatal(http.ListenAndServe(":8088", nil))

}


func handleWs(ws *websocket.Conn)  {

	var m sync.Mutex;

	var wait sync.WaitGroup;

	wait.Add(1)


	ch := make(chan int)

	go addSocketToMap(ws, &m, ch)

	key := <- ch

	fmt.Println(key)

	go listenToWs(ws, key, &wait)

	fmt.Println("New socket added success");

	wait.Wait()

}

func addSocketToMap(ws *websocket.Conn, m *sync.Mutex, ch chan int)  {

	m.Lock()

	socketMap[socketCount] =  ws;
	ch <- socketCount
	socketCount++

	m.Unlock()

}


func listenToWs(ws *websocket.Conn, key int, wait *sync.WaitGroup)  {

	for  {

		var reply string;

		websocket.Message.Receive(ws, &reply)

		fmt.Println("got message ", reply)

		go broadCastMessageToAll(reply, key)

	}

	wait.Done()

}


func broadCastMessageToAll(message string, socketKey int)  {

	for key, ws := range socketMap {

		if key == socketKey {
			continue
		}
		websocket.Message.Send(ws, message)
	}

}