package flow

import (
	"fmt"
	"reflect"
)

//Chan Flow Chan
type Chan interface{}

//Context Flow Context
type Context map[string]interface{}

//Process Flow Chan
type Process interface {
	Process(input Chan, output Chan, ctx Context)
}

//Flow flujo
type Flow struct {
	outputChan Chan
	item       Process
	prev       *Flow
	next       *Flow
	err        error
}

//New New Flow
func New(item Process, outputChan Chan) *Flow {
	flow := Flow{
		prev:       nil,
		item:       item,
		outputChan: outputChan,
	}
	return &flow
}

//Start a Flow
func (f *Flow) Start(inputChan Chan, ctx Context) Chan {
	out := f.makeChannel()
	f.run(inputChan, out, ctx)

	next := f.next
	in := out
	for next != nil {
		out = next.makeChannel()
		next.run(in, out, ctx)
		next = next.next
		in = out
	}
	return out
}

//Add Add
func (f *Flow) Add(item Process, outputChan Chan) *Flow {
	if f.next != nil {
		return f.next.Add(item, outputChan)
	}
	flow := Flow{
		prev:       f,
		item:       item,
		outputChan: outputChan,
	}
	f.next = &flow
	return f
}

//Run Run
func (f *Flow) run(input, output Chan, ctx Context) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.first().err = fmt.Errorf("Flow: '%v.Process()' ->  %v ", reflect.TypeOf(f.item), r)
			}
		}()
		f.item.Process(input, output, ctx)
	}()
}

func (f *Flow) makeChannel() Chan {
	if f.outputChan == nil {
		return make(chan interface{})
	}
	cap := 0
	cType := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(f.outputChan))
	if reflect.TypeOf(f.outputChan).Kind() == reflect.Chan {
		cap = reflect.ValueOf(f.outputChan).Cap()
		cType = reflect.ChanOf(reflect.BothDir, reflect.TypeOf(f.outputChan).Elem())
	}
	return reflect.MakeChan(cType, cap).Interface()
}

func (f *Flow) first() *Flow {
	if f.prev != nil {
		return f.prev.first()
	}
	return f
}

//Err flow error
func (f *Flow) Err() error {
	return f.first().err
}
