package flow_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/coc1961/flow/pkg/flow"
)

func TestFlow_run(t *testing.T) {
	Step1 := Step1{}
	Step2 := Step2{}
	Step3 := Step3{}

	f1 := flow.New(Step1, make(chan string, 0))
	f1.Add(Step2, make(chan int, 0))
	f1.Add(Step3, make(chan int, 0))

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			input := make(chan string, 1)
			outChan := f1.Start(input)
			for n := 0; n < 100; n++ {
				input <- "10"
				input <- "20"
				input <- "30"
				input <- "40"
			}
			close(input)

			out := outChan.(chan int)
			res := <-out
			fmt.Println(res)
			if res != 10000 {
				t.Error("Flow Error")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

type Step1 struct {
}

func (p Step1) Process(input flow.Chan, output flow.Chan) {
	fmt.Println("Start Step1")

	out := output.(chan string)
	in := input.(chan string)

	for {
		str, ok := <-in
		if ok {
			out <- str
		} else {
			break
		}
	}

	close(out)

	fmt.Println("End Step1")
}

type Step2 struct {
}

func (p Step2) Process(input flow.Chan, output flow.Chan) {
	fmt.Println("Start Step2")
	in := input.(chan string)
	ou := output.(chan int)

	for {
		str, ok := <-in
		if ok {
			i, _ := strconv.Atoi(str)
			ou <- i
		} else {
			break
		}
	}

	close(ou)

	fmt.Println("End Step2")
}

type Step3 struct {
}

func (p Step3) Process(input flow.Chan, output flow.Chan) {
	fmt.Println("Start Step3")
	in := input.(chan int)
	ou := output.(chan int)

	tot := 0

	for {
		num, ok := <-in
		if ok {
			tot += num
		} else {
			break
		}
	}

	ou <- tot
	close(ou)

	fmt.Println("End Step3")
}
