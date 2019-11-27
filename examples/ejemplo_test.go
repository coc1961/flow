package examples

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/coc1961/flow/pkg/flow"
)

func Test_Flow(t *testing.T) {
	//Cantidad de hilos de ejecución
	elements := 6

	//Canal de entrada
	input := make(chan []int)

	// Armo el Flujo
	// Parser: parte el array de int en partes
	// Summarizer: toma los arrays de int crea hilos y cada hilo suma una parte del total
	// Joiner: toma los totales que retorna Summarizer y retorna un solo total
	f1 := flow.New(Parser{Elements: elements}, input, make(chan [][]int, 0)).
		Add(Summarizer{}, make(chan int, 0)).
		Add(Joiner{}, make(chan int, 0))

	//Preparo Entrada, Lleno Array
	iArray := make([]int, 0, 10000000*elements)
	for i := 1; i <= cap(iArray); i++ {
		iArray = append(iArray, rand.Int())
	}

	//Ejecuto
	start := time.Now()
	input <- iArray

	//Salida
	out := f1.Out().(chan int)
	total := <-out
	fmt.Println("Con Flow:", time.Since(start))

	//Controlo resultado
	start = time.Now()
	total1 := 0
	for _, v := range iArray {
		total1 += v
	}
	fmt.Println("Sin Flow:", time.Since(start))

	fmt.Println("Status:", map[bool]string{true: "Ok", false: "Error"}[total == total1], ", Len:", len(iArray))

	if total != total1 {
		t.Errorf("Total error want %v, got %v", total1, total)
	}
}
