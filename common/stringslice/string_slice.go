package stringslice

func Intersection(a []string, b []string) []string {
	var r []string
	for _, ae := range a {
		for _, be := range b {
			if ae == be {
				r = append(r, ae)
			}
		}
	}
	return r
}

func Shift(a []string, b string) []string {
	r := []string{b}
	return append(r, a...)
}

//element should be unique
//for example ["a","b","b"] this is not allow
//same for the ElementEqual Method
func RemoveDupSS(ss [][]string) [][]string {
	var rss [][]string
	var exist bool
	for i := 0; i < len(ss); i++ {
		exist = false
		for j := i + 1; j < len(ss); j++ {
			if ElementEqual(ss[i], ss[j]) {
				exist = true
				break
			}
		}
		if !exist {
			rss = append(rss, ss[i])
		}
	}
	return rss
}

func ElementEqual(as []string, bs []string) bool {
	if len(as) != len(bs) {
		return false
	}
	for _, a := range as {
		if Contains(bs, a) == false {
			return false
		}
	}
	return true
}

func Contains(ss []string, s string) bool {
	for _, e := range ss {
		if e == s {
			return true
		}
	}
	return false
}

func Reverse(ss []string) []string {
	var rs = make([]string, len(ss))
	for i := range ss {
		rs[len(ss)-i-1] = ss[i]
	}
	return rs
}

func UniqueAppend(ss []string, cs ...string) []string {
	var exist bool
	for _, s := range cs {
		exist = false
		for _, e := range ss {
			if e == s {
				exist = true
				break
			}

		}
		if exist == false {
			ss = append(ss, s)
		}
	}

	return ss
}
