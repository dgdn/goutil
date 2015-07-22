package common

import (
	"strings"
)

type CVSStringSS []string

func (cs *CVSStringSS) ScanP(str string) error {
	var ss []string
	for _, s := range strings.Split(str, ",") {
		ss = append(ss, s)
	}
	*cs = ss
	return nil
}
