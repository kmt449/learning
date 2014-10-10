/* section 5
 * Producer-Consumer
 */

// -*- Encoding: UTF-8 -*-

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Producer struct {
	id    int
	outch chan int
}

func producerThread(pr *Producer) {
	for {
		pr.outch <- pr.id
		pr.id++
	}
}

type Maker struct {
	name  string
	pr    *Producer
	table chan string
	seed  int64
}

func MakerThread(mk *Maker) {
	rd := rand.New(rand.NewSource(mk.seed))
	for {
		mk.table <- fmt.Sprintf("[ Cake No. %d by %s ]", <-mk.pr.outch, mk.name)
		time.Sleep(time.Duration(rd.Intn(1000)) * time.Millisecond)
	}
}

type Eater struct {
	name  string
	table chan string
	seed  int64
}

func EaterThread(et *Eater) {
	rd := rand.New(rand.NewSource(et.seed))
	for {
		fmt.Printf("%s takes %s\n", et.name, <-et.table)
		time.Sleep(time.Duration(rd.Intn(1000)) * time.Millisecond)
	}
}

func main() {
	fmt.Println("hello")

	table := make(chan string, 3)

	pr := &Producer{0, make(chan int)}
	go producerThread(pr)

	go MakerThread(&Maker{"MakerThread-1", pr, table, 31415})
	go MakerThread(&Maker{"MakerThread-2", pr, table, 92653})
	go MakerThread(&Maker{"MakerThread-3", pr, table, 58979})
	go EaterThread(&Eater{"EaterThread-1", table, 32384})
	go EaterThread(&Eater{"EaterThread-2", table, 62643})
	go EaterThread(&Eater{"EaterThread-3", table, 38327})

	time.Sleep(5 * time.Second)
}
