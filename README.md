# Flow

Simple framework to design workflows based on processes connected through channels.

> Each process has an input channel and an output channel,the output channel is the input channel of the next process.

> Each process is responsible for closing the output channel when no more information will be sent.

>  Each process will process the data received from its input channel until the end of the data.

> Each created process can be reused in other flows

Creating a Process:

```go
type Summarizer struct {
}

func (p Summarizer) Process(input flow.Chan, output flow.Chan) {
    //Cast input and output channels
    in := input.(chan int)
    out := output.(chan int)

    //I add the input values ​​and accumulate the total
    tot:=0
    for {
        num, ok := <-in
        if !ok {
            break
        }
        tot += nuk
    }

    //Send the total to the output channel
    out <- tot

    //I close the output channel, no more data to send
    close(out)

}
```

Creating a Flow

```go
    // I create the flow with an instance of Process.
    // The input channel of the flow.
    // The output channel of the flow
    fl := flow.New(Summarizer{}, make(chan int, 0))

    // I create the flow data input channel
    input := make(chan string, 1)

    //Start the flow
    outChan := fl.Start(input)
```

Using the Flow

```go
    //Send data to the flow
    input <- 10
    input <- 20
    input <- 30
    input <- 40

    //No more data to send, I close the channel
    close(input)

    //Cast the flow outflow and get the result
    out := outChan.(chan int)
    total := <-out
    fmt.Println(total)

    //Printed Result is
    100
```

Adding Processes to the Flow

```go

    // I create the flow with an instance of Process.
    // The input channel of the flow.
    // The output channel of the flow
    fl := flow.New(Summarizer{}, input, make(chan int, 0))

    // I add a new process to the flow ToString
    // This process has as input channel the output channel of Summarizer
    // The output channel of the ToString process is of type string
    // This process takes the Summarize total and transforms it into a string
    fl.Add(ToString{}, make(chan string, 0))

    // I create the flow data input channel
    input := make(chan string, 1)

   //Start the flow
    outChan := fl.Start(input)

    // I execute the flow and now the total 100 number is transformed into a string as shown below

    //Cast the flow outflow and get the result
    out := outChan.(chan string)
    total := <-out
    fmt.Println(total)

    //Printed Result is
    "The Total is : 100"

```
