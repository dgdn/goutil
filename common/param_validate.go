package common

import (
	"fmt"
)

func ValidateParam(vs ...interface{}) error {
	for _, v := range vs {
		if v != "" {
			return fmt.Errorf("%v", v)
		}
	}
	return nil
}

//deprecated
func MustBePositive(v int64, tips ...interface{}) string {
	if v <= 0 {
		return fmt.Sprintf("%v", tips...)
	}
	return ""
}

//deprecated
func Int64In(v int64, vs []int64, tips ...interface{}) string {
	for _, e := range vs {
		if e == v {
			return ""
		}
	}
	return fmt.Sprintf("[%v] is not in %v", v, vs)
}

//deprecated
func StringIn(v string, vs []string, tips ...interface{}) string {
	for _, e := range vs {
		if e == v {
			return ""
		}
	}
	return fmt.Sprintf("[%v] is not in %v", v, vs)
}

func VP_MustBePositive(v int64, tips ...interface{}) string {
	if v <= 0 {
		return fmt.Sprintf("%v", tips...)
	}
	return ""
}

func VP_Int64In(v int64, vs []int64, tips ...interface{}) string {
	for _, e := range vs {
		if e == v {
			return ""
		}
	}
	return fmt.Sprintf("[%v] is not in %v", v, vs)
}

func VP_StringIn(v string, vs []string, tips ...interface{}) string {
	for _, e := range vs {
		if e == v {
			return ""
		}
	}
	return fmt.Sprintf("[%v] is not in %v", v, vs)
}

func VP_StringRange(v string, from int, to int, tips ...interface{}) string {
	if len(v) >= from && len(v) <= to {
		return ""
	}
	return fmt.Sprintf("[%v] len must satisfy %v<= len <=%v", tips, from, to)
}
