package unbuf

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Player xxx
type Player struct {
	User1 string
	User2 string
}

// NewPlayer  create a instance, Player is object
func NewPlayer(user1, user2 string) *Player {
	return &Player{
		User1: user1,
		User2: user2,
	}
}

func (p *Player) GoPlayer() {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(2)
	go player(p.User1, ch, &wg)
	go player(p.User2, ch, &wg)
	ch <- 1
	wg.Wait()
}

//player  two players
func player(name string, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		ball, ok := <-ch
		if !ok {
			fmt.Printf("%s won!!!\n", name)
			break
		}
		n := rand.Intn(100)
		if n%15 == 0 {
			fmt.Printf("%s miss,the rand number is %d\n", name, n)
			close(ch)
			break
		}
		fmt.Printf("Player %s hit the ball %d,the rand number is %d\n", name, ball, n)
		ball++
		ch <- ball
	}
}

func Runner() {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)
	go run(ch, &wg)
	ch <- 1
	wg.Wait()
}

func run(ch chan int, wg *sync.WaitGroup) {
	var newRunner int
	runner := <-ch
	fmt.Printf("runner %d running with Baton\n", runner)
	if runner != 4 {
		newRunner = runner + 1
		fmt.Printf("runner %d to the line\n", runner)
		go run(ch, wg)
	}
	// rand sleep time
	n := rand.Intn(100)
	time.Sleep(time.Duration(n) * time.Millisecond)
	if runner == 4 {
		fmt.Printf("runner %d finish,Race over\n", runner)
		wg.Done()
		return
	}
	fmt.Printf("runner %d exchange with runner %d\n", runner, newRunner)
	ch <- newRunner
}