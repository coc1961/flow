package flow

import (
	"reflect"
)

//Chan Flow Chan
type Chan interface{}

//Process Flow Chan
type Process interface {
	Process(input Chan, output Chan)
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
func (f *Flow) Start(inputChan Chan) Chan {
	if f.prev != nil {
		f.prev.Start(inputChan)
	}
	out := f.makeChannel()
	f.run(inputChan, out)

	next := f.next
	in := out
	for next != nil {
		out = next.makeChannel()
		next.run(in, out)
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
func (f *Flow) run(input, output Chan) {
	go f.item.Process(input, output)
}

func (f *Flow) makeChannel() interface{} {
	if f.outputChan == nil {
		return nil
	}
	cType := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(f.outputChan))
	if reflect.TypeOf(f.outputChan).Kind() == reflect.Chan {
		cType = reflect.ChanOf(reflect.BothDir, reflect.TypeOf(f.outputChan).Elem())
	}
	return reflect.MakeChan(cType, 0).Interface()
}
