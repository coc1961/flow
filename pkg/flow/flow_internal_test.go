package flow

import (
	"reflect"
	"testing"
)

func TestFlow_makeChannel(t *testing.T) {
	type fields struct {
		outputChan Chan
		item       Process
		prev       *Flow
		next       *Flow
	}
	tests := []struct {
		name   string
		fields fields
		want   Chan
	}{
		{
			name: "Nil chan",
			fields: fields{
				outputChan: nil,
			},
			want: make(chan interface{}, 0),
		},
		{
			name: "Int chan",
			fields: fields{
				outputChan: int(0),
			},
			want: make(chan int, 0),
		},
		{
			name: "Int chan",
			fields: fields{
				outputChan: make(chan int, 10),
			},
			want: make(chan int, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flow{
				outputChan: tt.fields.outputChan,
				item:       tt.fields.item,
				prev:       tt.fields.prev,
				next:       tt.fields.next,
			}

			got := f.makeChannel()
			v1 := reflect.ValueOf(got).Type()
			v2 := reflect.ValueOf(tt.want).Type()
			if !reflect.DeepEqual(v1, v2) {
				t.Errorf("Flow.makeChannel() = %v, want %v", v1, v2)
			}
		})
	}
}
