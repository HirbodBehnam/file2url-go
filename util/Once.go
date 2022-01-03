package util

// SingleThreadOnce is something like sync.Once but not thread safe
// It is very good if you want to run a "defer" statement once and early
type SingleThreadOnce struct {
	done bool
}

// Do run a function only if no other function was run before
func (o *SingleThreadOnce) Do(f func()) {
	if !o.done {
		o.done = true
		f()
	}
}
