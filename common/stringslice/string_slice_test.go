package stringslice

import (
	"reflect"
	"testing"
)

func TestIntersection(t *testing.T) {
	r := Intersection([]string{"a", "b", "c"}, []string{"b", "c", "d"})
	if len(r) != 2 {
		t.Errorf("expect %v but got %v", 2, len(r))
	}
	if r[0] != "b" && r[1] != "c" {
		t.Errorf("expect %v got %v", []string{"a", "b"}, r)
	}
	r1 := Intersection([]string{"a"}, []string{"b"})
	if len(r1) != 0 {
		t.Errorf("expect %v got %v", 0, len(r1))
	}
}

func TestRemoveDupSS(t *testing.T) {
	r := RemoveDupSS([][]string{
		[]string{"a", "b"},
		[]string{"b", "a"},
		[]string{"a", "b", "c"},
	})
	if len(r) != 2 {
		t.Errorf("expect %v but %v", 2, len(r))
	}
	if reflect.DeepEqual(r[0], []string{"b", "a"}) == false {
		t.Errorf("exect %v but %v", r[0], []string{"b", "a"})
	}
	if reflect.DeepEqual(r[1], []string{"a", "b", "c"}) == false {
		t.Errorf("exect %v but %v", r[0], []string{"a", "b", "c"})
	}
}

func TestUniqueAppend(t *testing.T) {
	s1 := UniqueAppend([]string{"a", "b"}, "a")
	if len(s1) != 2 {
		t.Errorf("expect %v but %v", 2, len(s1))
	}
	if s1[0] != "a" && s1[1] != "b" {
		t.Errorf("expect %v but %v", []string{"a", "b"}, s1)
	}

	s2 := UniqueAppend([]string{"a", "b"}, "c")
	if len(s2) != 3 {
		t.Errorf("expect %v but %v", 3, len(s2))
	}
	if s2[0] != "a" && s2[1] != "b" && s2[2] != "c" {
		t.Errorf("expect %v but %v", []string{"a", "b", "c"}, s2)
	}
}

func TestReverse(t *testing.T) {
	s1 := Reverse([]string{"a", "b"})
	if len(s1) != 2 {
		t.Errorf("expect %v got %v", 2, len(s1))
	}
	if s1[0] != "b" && s1[1] != "c" {
		t.Errorf("expect %v got %v", []string{"b", "a"}, s1)
	}
}
