package optional

// https://github.com/markphelps/optional MIT License
import (
	"encoding/json"
)

func New[T any](vs ...T) Option[T] {
	o := Option[T]{
		set: true,
	}
	if len(vs) > 0 {
		o.value = &vs[0]
	}

	return o
}

type Option[T any] struct {
	value *T
	set   bool
}

func (t *Option[T]) IsSet() bool {
	return t.set
}

func (t *Option[T]) Set(v T) {
	t.value = &v
	t.set = true
}

func (t Option[T]) Get() (T, bool) {
	if !t.IsPresent() {
		var zero T
		return zero, false
	}

	return *t.value, true
}
func (t Option[T]) MustGet() T {
	if !t.IsPresent() {
		panic("value is not present")
	}

	return *t.value
}

func (t Option[T]) IsPresent() bool {
	return t.value != nil
}

func (t Option[T]) OrElse(v T) T {
	if t.IsPresent() {
		return *t.value
	}
	return v
}

func (t Option[T]) If(fn func(T)) Option[T] {
	if t.IsPresent() {
		fn(*t.value)
	}
	return t
}
func (t Option[T]) MarshalJSON() ([]byte, error) {
	if t.IsPresent() {
		return json.Marshal(t.value)
	}
	return json.Marshal(nil)
}

func (t *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.value = nil
		t.set = true
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	t.value = &value
	t.set = true
	return nil
}
