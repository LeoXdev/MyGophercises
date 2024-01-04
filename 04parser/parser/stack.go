package parser

import "errors"

// stack type is a simple implementation of a stack,
// destined to check opening and closing html tags.
//
// With each opening html tag encountered, a new link object will open to be filled (push).
// And with each closing, a link object will be popped from the stack and moved
// to a links slice.
type stack []link

func NewStack() *stack {
	return new(stack)
}
func (s *stack) IsEmpty() bool {
	return len(*s) == 0
}
func (s *stack) Push(newLink link) {
	(*s) = append((*s), newLink)
}
func (s *stack) Pop() (link, error) {
	if (*s).IsEmpty() {
		return link{}, errors.New("stack is empty")
	}
	r := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]

	return r, nil
}
func (s *stack) Peek() (link, error) {
	if (*s).IsEmpty() {
		return link{}, errors.New("stack is empty")
	}
	return (*s)[len(*s)-1], nil
}