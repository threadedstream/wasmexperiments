package queue

type Queue interface {
	Push(val interface{})
	Peek() interface{}
	Pop()
	Empty() bool
}

type DummyQueue struct{}

func (dq DummyQueue) Push(_ interface{}) {}
func (dq DummyQueue) Peek() interface{}  { return nil }
func (dq DummyQueue) Pop()               {}
func (dq DummyQueue) Empty() bool        { return true }
