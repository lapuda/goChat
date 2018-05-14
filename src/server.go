package main

import (
	"fmt"
	"net"
	//"io"
	"os"
	"flag"
	"log"
	"sync/atomic"
)
var host = flag.String("host","","host")
var port = flag.String("port","3333","port")


type Connection struct {
	id   uint64
	conn net.Conn
}

var ID uint64  = 100012

var client = make([]Connection,1024)
func main(){
	flag.Parse()
	var l net.Listener
	var err error
	l,err = net.Listen("tcp",*host+":"+*port)
	if err != nil {
		fmt.Println("Error listening:",err)
		os.Exit(0)
	}
	defer l.Close()
	fmt.Println("Listen on " + *host + ":" + *port)
	for {
		conn,err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:",err)
			os.Exit(1)
		}
		fmt.Printf("收到消息: %s -> %s\n",conn.RemoteAddr(),conn.LocalAddr())
		var title = []byte("欢迎您:"+conn.RemoteAddr().String())
		conn.Write(title)
		ID++;
		var connection = Connection{ID,conn}
		client = append(client, connection)
		go handerRequest(conn,connection,conn.RemoteAddr().String())
	}
}

func GetIncreaseID(ID *uint64) uint64 {
	var n, v uint64
	for {
		v = atomic.LoadUint64(ID)
		n = v + 1
		if atomic.CompareAndSwapUint64(ID, v, n) {
			break
		}
	}
	return n
}

func handerRequest(conn net.Conn,connection Connection,name string)  {
	defer closeAll()
	buf := make([]byte,1024)
	for {
		//io.Copy(conn,conn)
		lengh,err := conn.Read(buf);
		if (err != nil) {
			log.Println(err)
			conn.Close()
			remove(connection.id)
			return
		}
		message := name + ":"+string(buf[0:lengh])
        sendAll(connection,message)
		log.Println(string(buf[0:lengh]))
	}
}

func remove(id uint64)  {
	for index,conn :=  range client {
		if (conn.id == id) {
			client = append(client[:index],client[index+1:]...)
		}
	}
}

func sendAll(conns Connection,message string)  {
	log.Println("收到:",conns.id)
	for _,conn :=  range client {
		if (conn.conn != nil && conns.id != conn.id) {
			conn.conn.Write([]byte(message))
			log.Println("发送给:" , conn.id)
		}
	}
}

func closeAll()  {
	for _,conn :=  range client {
		if (conn.conn != nil) {
			conn.conn.Close()
		}
	}
}