package stack

import (
	"errors"
	"testing"
)

func TestNewStackIsEmpty(t *testing.T) {
	b := NewStack[bool]()
	if !b.IsEmpty() {
		t.Errorf("%v isEmpty() returned %v; want true", b, b.IsEmpty())
	}
}

func TestNewStackString(t *testing.T) {
	b := NewStack[bool]()

	if b.String() != "[]" {
		t.Errorf("%v String() returned %v; want []", b, b.IsEmpty())
	}
}

func TestNewStackPopOrElse(t *testing.T) {
	e := errors.New("error")
	s := NewStack[string]().PopOrElse(func() interface{} { return e })

	if e != s {
		t.Errorf("NewStack[string]().PopOrElse(func() interface{} { return e }) returned %v; want %v", s, e)
	}
}

func TestNewStackPush(t *testing.T) {
	s := NewStack[string]().Push("Hello")

	if s.TopOrElse(func() interface{} {
		return errors.New("TopOrElse(func() interface{} {return error}) failed on a stack with 1 element")
	}) != "Hello" {
		t.Errorf("TopOrElse(func() interface{} {return error}) failed on a stack with 1 element")
	}
}

func TestStackPush(t *testing.T) {
	s := NewStack[string]().Push("Hello")

	s = s.Push("World")

	e := s.TopOrElse(func() interface{} {
		return errors.New("TopOrElse(func() interface{} {return error}) failed on a stack with 1 element")
	})

	if e != "World" {
		t.Errorf("TopOrElse(func() interface{} {return error}) returned %v; want World", e)
	}

	if s.String() != "[WorldHello]" {
		t.Errorf("[WorldHello].String() returned %v; want [WorldHello]", s.String())
	}
}

func TestStackPushCacheEntriesCorrect(t *testing.T) {
	s := NewStack[string]().Push("Hello").Push("World")

	if s.cache["World"].String() != "[Hello]" {
		t.Errorf("s.cache[\"World\"].String() returned %v; want [Hello]", s.String())
	}

	if s.cache["World"].cache["Hello"].String() != "[]" {
		t.Errorf("s.cache[\"World\"].cache[\"Hello\"].String() returned %v; want []", s.String())
	}

	if s.cache["Hello"].String() != "" {
		t.Errorf("s.cache[\"Hello\"] returned %v; want ", s.String())
	}

	if s.cache["World"].cache["World"].String() != "" {
		t.Errorf("s.cache[\"World\"].cache[\"World\"] returned %v; want ", s.String())
	}
}

func TestStackPopOrElse(t *testing.T) {
	e := errors.New("PopOrElse(func() interface{} {return error}) failed on a stack with 1 element")
	s := NewStack[string]().Push("Hello").Push("World")
	s2 := s.PopOrElse(func() interface{} {
		return e
	})

	if s2 == e {
		t.Errorf("PopOrElse(func() interface{} {return error}) failed on a stack with 1 element")
	}

	stackResult := (s2).(stack[string])

	if stackResult.String() != "[Hello]" {
		t.Errorf("[Hello].String() returned %v; want [Hello]", s.String())
	}
}

func TestPushPopPush(t *testing.T) {
	e := errors.New("should not throw")
	s := (NewStack[string]().
		Push("Hello").
		Push("World").
		PopOrElse(func() interface{} { return e })).(stack[string]).
		Push("Daisy")

	if s.String() != "[DaisyHello]" {
		t.Errorf("Got %v; expected [DaisyHello]", s.String())
	}

	elem := s.TopOrElse(func() interface{} { return e })

	if elem != "Daisy" {
		t.Errorf("Got %v; expected [Daisy]", s.String())
	}
}

func TestPushPopPushSameKey(t *testing.T) {
	e := errors.New("should not throw")
	s := (NewStack[string]().
		Push("Hello").
		Push("World").
		PopOrElse(func() interface{} { return e })).(stack[string]).
		Push("World")

	if s.String() != "[WorldHello]" {
		t.Errorf("Got %v; expected [WorldHello]", s.String())
	}

	elem := s.TopOrElse(func() interface{} { return e })

	if elem != "World" {
		t.Errorf("Got %v; expected [World]", s.String())
	}
}

func TestPushSameKeyTwice(t *testing.T) {
	e := errors.New("should not throw")
	s := NewStack[string]().
		Push("Hello").
		Push("World").
		Push("World")

	if s.String() != "[WorldWorldHello]" {
		t.Errorf("Got %v; expected [WorldWorldHello]", s.String())
	}

	// get the top elem
	elem := s.TopOrElse(func() interface{} { return e })
	if elem != "World" {
		t.Errorf("Got %v; expected World", s.String())
	}

	// pop an elem
	s = (s.PopOrElse(func() interface{} { return e })).(stack[string])
	if s.String() != "[WorldHello]" {
		t.Errorf("Got %v; expected [WorldHello]", s.String())
	}
	// get the top elem
	elem = s.TopOrElse(func() interface{} { return e })
	if elem != "World" {
		t.Errorf("Got %v; expected World", s.String())
	}

	// pop an elem
	s = (s.PopOrElse(func() interface{} { return e })).(stack[string])
	if s.String() != "[Hello]" {
		t.Errorf("Got %v; expected [Hello]", s.String())
	}
	elem = s.TopOrElse(func() interface{} { return e })
	if elem != "Hello" {
		t.Errorf("Got %v; expected Hello", s.String())
	}
}

func TestStackPopOrElseCacheEntriesCorrect(t *testing.T) {
	e := errors.New("PopOrElse(func() interface{} {return error}) failed on a stack with 1 element")
	s := NewStack[string]().Push("Hello").Push("World")
	s2 := s.PopOrElse(func() interface{} {
		return e
	})

	if s2 == e {
		t.Errorf("PopOrElse(func() interface{} {return error}) failed on a stack with 1 element")
	}

	stack := (s2).(stack[string])

	if stack.cache["Hello"].String() != "[]" {
		t.Errorf("s.cache[\"Hello\"].String() returned %v; want [Hello]", stack.cache["Hello"].String())
	}

	if stack.cache["World"].String() != "" {
		t.Errorf("s.cache[\"World\"].String() returned %v; want []", stack.cache["World"].String())
	}
}
