package pipelines

import "testing"

func TestPipelines(t *testing.T) {

	main := func() int {
		// Set up a done channel that's shared by the whole pipeline,
		// and close that channel when this pipeline exits, as a signal
		// for all the goroutines we started to exit.
		done := make(chan struct{})
		defer close(done)

		in := gen(done, 2, 3)

		// Distribute the sq work across two goroutines that both read from in.
		c1 := sq(done, in)
		c2 := sq(done, in)

		// Consume the first value from output.
		out := merge(done, c1, c2)
		res := <-out // 4 or 9
		return res

		// done will be closed by the deferred call.
	}

	if res := main(); res != 4 && res != 9 {
		t.Fail()
	}
}
