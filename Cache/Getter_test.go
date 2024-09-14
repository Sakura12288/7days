package Cache

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var ff Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("我很好")
	if v, _ := ff.Get("我很好"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
