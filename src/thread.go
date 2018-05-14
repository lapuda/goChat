package main

import (
	"fmt"
    "sync"
	"math/rand"
	"time"
)
const MAX = 10

type Routine struct {
	id    int
	xtime time.Time
}

type Result struct {
    showTime time.Time
    id int
}

var complete chan int = make(chan int)
var waitGroup sync.WaitGroup

//swap()
func showTime(info Routine,results chan<- *Result) {
    times := 1// rand.Intn(10)
	for i := 1; i < times; i++ {
		fmt.Println("ID:", info.id, "时间:", info.xtime)
		time.Sleep(time.Second)
	}
    var result *Result = new(Result)
    result.showTime = info.xtime
    result.id       = info.id
    results <- result
    waitGroup.Done()
}


//func showResult()

func main() {
	// goroutine showTime(Routine{rand.Intn(10000),time.Now()})
    defer fmt.Println("所有线程都执行完毕.....")
    results := make(chan *Result)
    waitGroup.Add(2*MAX)
	for i:= 1;i<= MAX;i++ {
        go showTime(Routine{rand.Intn(10000), time.Now()},results)
        go showTime(Routine{rand.Intn(10000), time.Now()},results)
    }
    
    go func() {
        // 等待所有的线程执行完毕
        waitGroup.Wait()
        close(results)
    }()
    for re := range results {
        fmt.Printf("时间%s,ID:%d\n",re.showTime,re.id)
    }
}
