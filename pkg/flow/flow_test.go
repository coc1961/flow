package flow_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/coc1961/flow/pkg/flow"
)

func TestFlow_run(t *testing.T) {
	paso1 := Paso1{}
	paso2 := Paso2{}
	paso3 := Paso3{}

	input := make(chan string, 1)

	//Creo los Flujos
	//Flujo1 out chan de string
	f1 := flow.New(paso1, input, make(chan string, 0))
	//Flujo2 in chan string (ou de Flujo1) ou chan int
	f1.Add(paso2, make(chan int, 0))

	//Flujo3 in chan int (ou de Flujo2) ou chan int
	f1.Add(paso3, make(chan int, 0))

	input <- "10"
	input <- "20"
	input <- "30"
	input <- "40"
	close(input)

	//Tomo el out del ultimo flujo
	out := f1.Out().(chan int)
	res := <-out
	fmt.Println(res)
}

type Paso1 struct {
}

func (p Paso1) Process(input flow.Chan, output flow.Chan) {
	// Prime Flujo, input es nil
	fmt.Println("Start Paso1")

	//Envío un string al canal de salida
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

	fmt.Println("End Paso1")
}

type Paso2 struct {
}

func (p Paso2) Process(input flow.Chan, output flow.Chan) {
	// Segundo flujo, recibe los datos en input del output del primer flujo
	fmt.Println("Start Paso2")
	in := input.(chan string)
	ou := output.(chan int)

	for {
		str, ok := <-in
		if ok {

			//Convierto el string en numérico y genero la salida en out
			i, _ := strconv.Atoi(str)
			ou <- i
		} else {
			break
		}
	}

	close(ou)

	fmt.Println("End Paso2")
}

type Paso3 struct {
}

func (p Paso3) Process(input flow.Chan, output flow.Chan) {
	// Segundo flujo, recibe los datos en input del output del primer flujo
	fmt.Println("Start Paso3")
	in := input.(chan int)
	ou := output.(chan int)

	tot := 0

	//Recibo input out de Paso2 y sumarizo
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

	fmt.Println("End Paso3")
}
