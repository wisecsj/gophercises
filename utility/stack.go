
type data interface{}

type stack struct {
	arr []data
}

func (s *stack) Push(val data) {
	s.arr = append(s.arr, val)
}
func (s *stack) Pop() data {
	l := len(s.arr)
	ret := s.arr[l-1]
	s.arr = s.arr[:l-1]
	return ret
}
func (s *stack) isEmpty() bool {
	return len(s.arr) == 0
}