package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Package struct {
	Type  string            `json:"type"`
	Data  map[string]string `json:"data"`
	Key   string            `json:"key"`
	Value string            `json:"value"`
}

var mutex = &sync.Mutex{}
var list map[string]string

func main() {
	list = make(map[string]string)
	fmt.Println("Launching server...")

	target := flag.String("d", "", "target peer to dial")
	flag.Parse()

	ln, _ := net.Listen("tcp", ":0")
	fmt.Println("Using port:", ln.Addr().(*net.TCPAddr).Port)
	go ListenConn(ln)

	if *target != "" {
		conn, _ := net.Dial("tcp", *target)

		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

		go writeData(rw)
		go readData(rw)

	}
	select {}

}

func ListenConn(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		if err != nil {
			fmt.Println("Broken pipeline")
		}
		go handleStream(conn)
	}
}

func handleStream(s net.Conn) {

	log.Println("Got a new stream!")

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')

		if err != nil {
			//To shutdown nodes if one pipe is broken
			//log.Fatal(err)
		}

		if str == "" {
			return
		}

		if str != "\n" {

			p := Package{}
			err := json.Unmarshal([]byte(str), &p)
			if err != nil {
				log.Fatal(err)
			}

			mutex.Lock()

			switch p.Type {
			case "add":
				list[p.Key] = p.Value
				printChanges()
			case "exchange":
				if len(p.Data) > len(list) {
					list = p.Data
					printChanges()
				}
			case "get":
				if val, ok := list[p.Key]; ok {
					rw.Flush()
					rw.WriteString(val + "\n")
				} else {
					rw.Write([]byte("key not found\n"))
				}
			}
			rw.Flush()
			mutex.Unlock()
		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		p := Package{
			Type: "exchange",
			Data: list,
		}
		bytes, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
		}
		mutex.Unlock()

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))

		rw.Flush()
		mutex.Unlock()

	}
}

func printChanges() {
	bytes, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	//	// Green console color: 	\x1b[32m
	//	// Reset console color: 	\x1b[0m
	fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
}
