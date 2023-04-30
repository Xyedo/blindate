package optional

// https://github.com/markphelps/optional MIT License
import (
	"encoding/json"
)

type Option[T comparable] struct {
	value *T
	set   bool
}

func (t *Option[T]) JSONKeySent() bool {
	return t.set
}

func (t *Option[T]) Set(v T) {
	t.value = &v
}

func (t Option[T]) Get() (T, bool) {
	if !t.Present() {
		var zero T
		return zero, false
	}

	return *t.value, true
}
func (t Option[T]) MustGet() T {
	if !t.Present() {
		panic("value is not present")
	}

	return *t.value
}

func (t Option[T]) Present() bool {
	return t.value != nil
}

func (t Option[T]) OrElse(v T) T {
	if t.Present() {
		return *t.value
	}
	return v
}

func (t Option[T]) If(fn func(T)) Option[T] {
	if t.Present() {
		fn(*t.value)
	}
	return t
}
func (t Option[T]) MarshalJSON() ([]byte, error) {
	if t.Present() {
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
