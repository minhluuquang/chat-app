package trace

import (
	"fmt"
	"io"
)

// Tracer will help me to trace everything I do in this program
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct {
	out io.Writer
}

func (t *nilTracer) Trace(...interface{}) {}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// New return a new Tracer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// Off return a nilTracer (to turn off tracer)
func Off() Tracer {
	return &nilTracer{}
}
