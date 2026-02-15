package elevator

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Load values from a config file
//
//  Key-value pairs in the config file are assumed to be of the form:
//  "--key value"
//  Lines not starting in "--" are ignored.
//  Keys are *not* case-sensitive
//  Enum values are *not* case-sensitive
//
// arguments:
//  file:   Name of the file to load.
//  cases:  One or more instance of `con_val()` or `con_enum()`
//
// Example (Go-variant that stays close to the C-DSL):
//
//    type En int
//    const (En1 En = iota; En2; En3)
//
//    var i int
//    var s string
//    var en En
//
//    con_load("config.con",
//        con_val("integer", &i, "%d"),
//        con_val("greeting", &s, "%s"),
//        con_enum("enumeration", &en,
//            con_match("En1", En1),
//            con_match("En2", En2),
//            con_match("En3", En3),
//        ),
//    )
//

// --- Internal "case" interface (replaces the macro trickery) ---

type con_case interface {
	tryApply(key, val string)
}

// con_load(file, cases)
func con_load(file string, cases ...con_case) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Unable to open config file %s\n", file)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "--") {
			continue
		}

		// Mimic: sscanf(_line, "--%s %s", _key, _val);
		// -> key and value are whitespace-separated tokens (value can't contain spaces)
		fields := strings.Fields(line[2:])
		if len(fields) < 2 {
			continue
		}
		_key := fields[0]
		_val := fields[1]

		for _, c := range cases {
			c.tryApply(_key, _val)
		}
	}
}

// con_val(key, var, fmt)
type conValCase struct {
	key string
	ptr any
	fmt string
}

func con_val(key string, variablePtr any, fmtStr string) con_case {
	return conValCase{key: key, ptr: variablePtr, fmt: fmtStr}
}

func (c conValCase) tryApply(k, v string) {
	if !strings.EqualFold(k, c.key) {
		return
	}

	// Support common case where caller wants string "as is"
	// (C example used "%[^\n]" but their _val is from "%s", so no spaces anyway)
	if sp, ok := c.ptr.(*string); ok && (c.fmt == "%s" || c.fmt == "%[^\n]") {
		*sp = v
		return
	}

	_, _ = fmt.Sscanf(v, c.fmt, c.ptr)
}

// con_enum(key, var, match_cases)
type conEnumCase struct {
	key     string
	ptr     any
	matches []conEnumMatch
}

type conEnumMatch struct {
	name  string
	value any
}

func con_enum(key string, variablePtr any, matchCases ...conEnumMatch) con_case {
	return conEnumCase{key: key, ptr: variablePtr, matches: matchCases}
}

func con_match(id string, value any) conEnumMatch {
	return conEnumMatch{name: id, value: value}
}

func (c conEnumCase) tryApply(k, v string) {
	if !strings.EqualFold(k, c.key) {
		return
	}

	for _, m := range c.matches {
		if !strings.EqualFold(v, m.name) {
			continue
		}

		// Assign like C's: *var = _v;
		dst := reflect.ValueOf(c.ptr)
		if dst.Kind() != reflect.Pointer || dst.Elem().Kind() == reflect.Invalid {
			return
		}

		target := dst.Elem()
		src := reflect.ValueOf(m.value)

		// Allow assignment if types match or are convertible
		if src.Type().AssignableTo(target.Type()) {
			target.Set(src)
			return
		}
		if src.Type().ConvertibleTo(target.Type()) {
			target.Set(src.Convert(target.Type()))
			return
		}
		return
	}
}
