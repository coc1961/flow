package flow

//Chan Flow Chan
type Chan interface{}

//Item Flow Chan
type Item interface {
	Process(input Chan, output Chan)
}

//Flow flujo
type Flow struct {
	outputChan Chan
	item       Item
	prev       *Flow
	next       *Flow
}

//New New Flow
func New(item Item, inputChan, outputChan Chan) *Flow {
	flow := Flow{
		prev:       nil,
		item:       item,
		outputChan: outputChan,
	}
	flow.run(inputChan)
	return &flow
}

//Add Add
func (f *Flow) Add(item Item, outputChan Chan) *Flow {
	if f.next != nil {
		return f.next.Add(item, outputChan)
	}
	flow := Flow{
		prev:       f,
		item:       item,
		outputChan: outputChan,
	}
	f.next = &flow
	flow.run(nil)
	return f
}

//Run Run
func (f *Flow) run(input Chan) {
	if f.prev == nil {
		go f.item.Process(input, f.outputChan)
	} else {
		go f.item.Process(f.prev.outputChan, f.outputChan)
	}
}

//Out run
func (f *Flow) Out() Chan {
	if f.next != nil {
		return f.next.Out()
	}
	return f.outputChan
}
