package examples

import (
	"github.com/coc1961/flow/pkg/flow"
	"sync"
)

//Parser Parser
type Parser struct {
	Elements int
}

//Process Process (Implement Process)
func (p Parser) Process(input flow.Chan, output flow.Chan) {
	in := input.(chan []int)
	ou := output.(chan [][]int)

	//Divido el array en 4 subarray
	inputArray := <-in
	size := len(inputArray) / p.Elements

	//Creo un array de array de int con 4 elementos
	parsedArray := make([][]int, 0)
	for i := 0; i < len(inputArray); i += size {
		parsedArray = append(parsedArray, inputArray[i:i+size])
	}

	//Envío al proximo Flujo
	ou <- parsedArray
	close(ou)
}

//Summarizer Summarizer
type Summarizer struct {
}

//Process Process (Implement Process)
func (p Summarizer) Process(input flow.Chan, output flow.Chan) {
	in := input.(chan [][]int)
	ou := output.(chan int)

	parsedArray := <-in

	wg := sync.WaitGroup{}
	wg.Add(len(parsedArray))
	for ind := range parsedArray {
		go func(wg *sync.WaitGroup, ind int, out chan int) {
			tot := 0
			for _, v := range parsedArray[ind] {
				tot += v
			}
			out <- tot
			wg.Done()
		}(&wg, ind, ou)
	}
	wg.Wait()
	close(ou)
}

//Joiner Joiner
type Joiner struct {
}

//Process Process (Implement Process)
func (p Joiner) Process(input flow.Chan, output flow.Chan) {
	in := input.(chan int)
	ou := output.(chan int)

	tot := 0
	for {
		d, ok := <-in
		if !ok {
			break
		}
		tot += d
	}
	ou <- tot
	close(ou)
}
