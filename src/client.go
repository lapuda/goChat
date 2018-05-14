package main


import (
	"fmt"
	"net"
	"os"
	"bufio"
	"log"
	"time"
)
func checkError(err error) {
	if err != nil {
		fmt.Printf("Error : %s",err.Error())
		os.Exit(1)
	}
}

func main()  {
	conn,err := net.Dial("tcp","127.0.0.1:3333")
	checkError(err)
	defer conn.Close()
	go reciver(conn)
	send(conn)

}

func send(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		data,_,_ := reader.ReadLine()
		message := string(data)
		conn.Write([]byte(message))
	}
}

func reciver(conn net.Conn)  {
	buf := make([]byte,1024)
	for {
		length,err := conn.Read(buf)
		if (err != nil) {
			log.Println(err)
			os.Exit(1);
		}
		log.Println(string(buf[0:length])+" at "+time.Now().Format("2006-01-02 15:04:05"))
	}
}