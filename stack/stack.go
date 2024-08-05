package stack

type IStack[T comparable] interface {
	// push an element on the stack
	Push(elem T) stack[T]
	// get the element at the top of the stack, otherwise return
	// result from s2
	TopOrElse(s2 func() interface{}) interface{}
	// pop the element at the top of the stack, and return stack
	// otherwise return result from s2
	PopOrElse(s2 func() interface{}) interface{}
	// return a string representation of the stack
	String() string
	// returns true if no elements in the stack, false otherwise
	IsEmpty() bool
}

type stack[T comparable] struct {
	// collection of elements to following stacks
	// where T is type in the stack
	cache map[T]stack[T]
	// match function to be used on match call.
	//
	// changes via internal constructors to use
	// different functions depending on empty or
	// populated stack
	matchFunc func(onEmpty func() interface{}, onElem func(T, stack[T]) interface{}) interface{}
}

// get a new stack with an empty cache
// returns onEmpty when string and
func NewStack[T comparable]() stack[T] {
	return stack[T]{cache: make(map[T]stack[T]),
		matchFunc: func(onEmpty func() interface{}, onElem func(T, stack[T]) interface{}) interface{} {
			return onEmpty()
		}}
}

// returns new stack on element push
// previous stack cache becomes new stack cache
func newStackOnElem[T comparable](elem T, self stack[T]) stack[T] {
	return stack[T]{cache: make(map[T]stack[T]),
		matchFunc: func(onEmpty func() interface{}, onElem func(T, stack[T]) interface{}) interface{} {
			return onElem(elem, self)
		}}
}

// match function for top, pop and isEmpty operations
//
// Params:
//   - onEmpty() interface{}: Function to exec on empty stack. R type of object to return on error.
//   - onElem(T, stack[T]) interface{}: Function to exec on element in stack.
//     T type of elements in the stack, R type of elements to throw on error.
//     Returns an interface, which is actually the stack. We cannot return a type stack[T, R]
//     since we reuse the match methods for IsEmpty() and String() implementations
func match[T comparable](onEmpty func() interface{}, onElem func(T, stack[T]) interface{}, s stack[T]) interface{} {
	return s.matchFunc(onEmpty, onElem)
}

// match function for string
//
// Params:
//   - onEmpty() interface{}: Function to exec on empty stack. R type of object to return on error.
//   - onLast(T) interface{}: Function to execute when on the last element of the stack.
//   - onElem(T, stack[T]) interface{}: Function to exec on element in stack.
//     T type of elements in the stack, R type of elements to throw on error.
//     Returns an interface, which is actually the stack. We cannot return a type stack[T, R]
//     since we reuse the match methods for IsEmpty() and String() implementations
func match3[T comparable](onEmpty func() interface{}, onLast func(T) interface{}, onElem func(T, stack[T]) interface{}, s stack[T]) interface{} {
	return match(onEmpty, func(e T, t stack[T]) interface{} {
		if t.IsEmpty() {
			return onLast(e)
		}
		return onElem(e, t)
	}, s)
}

// internal push
//
// used to call internal function to create new stack
// with elem at top
func push[T comparable](elem T, s stack[T]) stack[T] {
	stack := newStackOnElem(elem, s)
	stack.cache[elem] = s
	return stack
}

// push a new element on the stack
func (s stack[T]) Push(elem T) stack[T] {
	return push(elem, s)
}

// get the element at the top of the stack, otherwise return
// result from onEmpty
func (s stack[T]) TopOrElse(onEmpty func() interface{}) interface{} {
	return match(onEmpty, func(e T, t stack[T]) interface{} {
		return e
	}, s)
}

// pop the element at the top of the stack, and return stack
// otherwise return result from onEmpty
func (s stack[T]) PopOrElse(onEmpty func() interface{}) interface{} {
	return match(onEmpty, func(e T, t stack[T]) interface{} {
		delete(t.cache, e)
		return t
	}, s)
}

func tostring[T comparable](s stack[T]) interface{} {
	return match3(func() interface{} { return "]" },
		func(e T) interface{} { str := any(e).(string); return str + "]" },
		func(e T, t stack[T]) interface{} {
			str := any(e).(string)
			ext := any(tostring(t)).(string)
			return str + ext
		},
		s)
}

// return a string representation of the stack
func (s stack[T]) String() string {
	var str = tostring(s)
	var toString = str.(string)

	return "[" + toString
}

func isempty[T comparable](s stack[T]) interface{} {
	return match(func() interface{} { return true }, func(e T, t stack[T]) interface{} { return false }, s)
}

// returns true if no elements in the stack, false otherwise
func (s stack[T]) IsEmpty() bool {
	var b = isempty(s)
	var toBool = b.(bool)

	return toBool
}
