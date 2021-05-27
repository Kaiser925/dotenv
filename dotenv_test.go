package dotenv

import (
	"os"
	"testing"
)

func TestLoadMap(t *testing.T) {
	envs := map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}
	err := LoadMap(envs)
	if err != nil {
		t.Fatal(err.Error())
	}

	for k, val := range envs {
		v := os.Getenv(k)
		if v != val {
			t.Fatalf("Excepted %s of %s, actual %s", val, k, v)
		}
	}
}

func TestParseKV(t *testing.T) {
	tests := []struct {
		line string
		k    string
		v    string
		err  error
	}{
		{
			line: "asdf",
			k:    "",
			v:    "",
			err:  errInvalidExpr,
		},
		{
			line: "A=a",
			k:    "A",
			v:    "a",
			err:  nil,
		},
		{
			line: "A = a",
			k:    "A",
			v:    "a",
			err:  nil,
		},
		{
			line: "A=\"a\"",
			k:    "A",
			v:    "a",
			err:  nil,
		},
		{
			line: "A='a'",
			k:    "A",
			v:    "a",
			err:  nil,
		},
		{
			line: "export A=\"a\"",
			k:    "A",
			v:    "a",
			err:  nil,
		},
		{
			line: "exportA=a",
			k:    "exportA",
			v:    "a",
			err:  nil,
		},
		{
			line: "export A=a # this is comment",
			k:    "A",
			v:    "a",
			err:  nil,
		},
		{
			line: "export A=\"\na \"",
			k:    "A",
			v:    "\na ",
			err:  nil,
		},
		{
			line: "export A B=a",
			k:    "",
			v:    "",
			err:  errKeyContainsSpace,
		},
		{
			line: "INVALID LINE",
			k:    "",
			v:    "",
			err:  errInvalidExpr,
		},
	}

	for _, c := range tests {
		k, v, err := parseKV(c.line)
		if err != c.err {
			t.Fatalf("excepted %v, actual %v", c.err, err)
		}

		if v != c.v {
			t.Fatalf("excepted %v, actual %v", c.v, v)
		}

		if k != c.k {
			t.Fatalf("excepted %v, actual %v", c.k, k)
		}
	}
}
