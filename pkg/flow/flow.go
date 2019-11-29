package flow

import (
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
	go f.item.Process(input, output, ctx)
}

func (f *Flow) makeChannel() Chan {
	if f.outputChan == nil {
		return make(chan interface{}, 0)
	}
	cType := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(f.outputChan))
	if reflect.TypeOf(f.outputChan).Kind() == reflect.Chan {
		cType = reflect.ChanOf(reflect.BothDir, reflect.TypeOf(f.outputChan).Elem())
	}
	return reflect.MakeChan(cType, 0).Interface()
}
