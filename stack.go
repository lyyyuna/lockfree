package lockfree

type Stack struct {
	top *stackNode
	len uint64
}

type stackNode struct {
	val  int
	next *stackNode
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(v int) {

}

func (s *Stack) Pop() int {
	return 0
}
