package queue

// type node[T comparable] struct {
// 	Data T
// 	next *node[T]
// }

// type list[T comparable] struct {
// 	start  *node[T]
// 	length int
// }

// type List[T comparable] interface {
// 	IsExist(data T) bool
// 	Pop(data T)
// 	Push(data T)
// 	Top() *node[T]
// 	Len() int
// }

// func NewList[T comparable]() List[T] {
// 	return &list[T]{
// 		start: nil,
// 	}
// }

// func (l *list[T]) IsExist(data T) bool {
// 	curr := l.start
// 	for curr != nil {
// 		if curr.Data == data {
// 			return true
// 		}
// 		curr = curr.next
// 	}
// 	return false
// }

// func (l *list[T]) Pop(data T) {
// 	if l.start == nil {
// 		return
// 	}

// 	l.length--

// 	if l.start.Data == data {
// 		l.start = l.start.next
// 		return
// 	}
// 	curr := l.start
// 	for curr.next != nil {
// 		if curr.next.Data == data {
// 			curr.next = curr.next.next
// 			return
// 		}
// 		curr = curr.next
// 	}
// }

// func (l *list[T]) Push(data T) {
// 	l.length++
// 	if l.start == nil {
// 		l.start = &node[T]{
// 			Data: data,
// 			next: nil,
// 		}
// 	} else {
// 		dummy := &node[T]{
// 			Data: data,
// 			next: l.start,
// 		}
// 		l.start = dummy
// 	}
// }

// func (l *list[T]) Top() *node[T] {
// 	return l.start
// }

// func (l *list[T]) Len() int {
// 	return l.lengthring
// }
