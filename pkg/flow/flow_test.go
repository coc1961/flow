package flow_test

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"

	"github.com/coc1961/flow/pkg/flow"
)

func TestFlow_run(t *testing.T) {
	Step1 := Step1{}
	Step2 := Step2{}
	Step3 := Step3{}
	Panic := Panic{}

	f1 := flow.New(Step1, make(chan string, 10))
	f1.Add(Step2, make(chan int, 10))
	f1.Add(Step3, make(chan int, 10))
	f1.Add(Panic, make(chan int, 10))

	wg := sync.WaitGroup{}

	for i := 1; i <= 1000; i++ {
		wg.Add(1)

		go func(ii int) {
			ctx := flow.Context{}

			ctx["counter"] = ii

			input := make(chan string, 1)
			outChan := f1.Start(input, ctx)
			for n := 0; n < 1000; n++ {
				input <- "10"
				input <- "20"
				input <- "30"
				input <- "40"
			}
			close(input)

			out := outChan.(chan int)
			res := <-out
			fmt.Println(res)
			if res != 100000 {
				t.Error("Flow Error")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println(f1.Err())
}

type Step1 struct {
}

func (p Step1) Process(input flow.Chan, output flow.Chan, ctx flow.Context) {
	fmt.Println("Start Step1", ctx["counter"])

	out := output.(chan string)
	in := input.(chan string)
	defer close(out)

	for {
		str, ok := <-in
		if ok {
			out <- str
		} else {
			break
		}
	}

	fmt.Println("End Step1", ctx["counter"])
}

type Step2 struct {
}

func (p Step2) Process(input flow.Chan, output flow.Chan, ctx flow.Context) {
	fmt.Println("Start Step2", ctx["counter"])
	in := input.(chan string)
	ou := output.(chan int)
	defer close(ou)

	for {
		str, ok := <-in
		if ok {
			i, _ := strconv.Atoi(str)
			ou <- i
		} else {
			break
		}
	}

	fmt.Println("End Step2", ctx["counter"])
}

type Step3 struct {
}

func (p Step3) Process(input flow.Chan, output flow.Chan, ctx flow.Context) {
	fmt.Println("Start Step3", ctx["counter"])
	in := input.(chan int)
	ou := output.(chan int)
	defer close(ou)

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

	fmt.Println("End Step3", ctx["counter"])
}

type Panic struct {
}

func (p Panic) Process(input flow.Chan, output flow.Chan, ctx flow.Context) {
	fmt.Println("Start Panic", ctx["counter"])
	in := input.(chan int)
	ou := output.(chan int)
	defer close(ou)
	for {
		num, ok := <-in
		if ok {
			ou <- num
		} else {
			break
		}
	}

	fmt.Println("End Panic", ctx["counter"])

	//Force Error to test Recovery
	var i *int
	*i++ // Nil force error
}
