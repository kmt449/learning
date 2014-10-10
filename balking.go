/* section 4
 * Balking
 */

// -*- Encoding: UTF-8 -*-

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Data struct {
	cond    *sync.Cond
	content string
	changed bool
}

func (dat *Data) init() {
	mutex := new(sync.Mutex)
	dat.cond = sync.NewCond(mutex)
	dat.changed = false
}

func (dat *Data) change(newContent string) {
	dat.cond.L.Lock()
	defer dat.cond.L.Unlock()
	dat.content = newContent
	dat.changed = true
	dat.cond.Signal()
}

func (dat *Data) save() {
	/* ここでロックしないと、期待の動作をしないことが確認できる */
	dat.cond.L.Lock()
	defer dat.cond.L.Unlock()
	if dat.changed {
		time.Sleep(1 * time.Second) /* バグを再現させるため */
		fmt.Printf("save: %s\n", dat.content)
		dat.changed = false
	}
}

type ChangerThread struct {
	name   string
	random *rand.Rand
	dat    *Data
}

func NewChangerThread(dat *Data, name string, seed int64) *ChangerThread {
	r := rand.New(rand.NewSource(seed))
	return &ChangerThread{name, r, dat}
}

func (ct *ChangerThread) Start(cf chan int) {
	for i := 0; i < 10; i++ {
		ct.dat.change(fmt.Sprintf("No.%d", i))
		time.Sleep(time.Duration(ct.random.Intn(1000)) * time.Millisecond)
		ct.dat.save()
	}
	cf <- 0
}

type ServerThread struct {
	name   string
	random *rand.Rand
	dat    *Data
}

func NewServerThread(dat *Data, name string, seed int64) *ServerThread {
	r := rand.New(rand.NewSource(seed))
	return &ServerThread{name, r, dat}
}

func (st *ServerThread) Start(cf chan int) {
	for i := 0; i < 10; i++ {
		st.dat.save()
		time.Sleep(time.Duration(st.random.Intn(1000)) * time.Millisecond)
	}
	cf <- 0
}

func main() {
	fmt.Println("hello")

	var dat Data
	dat.init()
	cf := make(chan int, 2)
	go NewServerThread(&dat, "serevr", 3141592).Start(cf)
	go NewChangerThread(&dat, "changer", 6535897).Start(cf)
	<-cf
	<-cf
}
