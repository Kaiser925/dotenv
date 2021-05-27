package dotenv

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

var errInvalidExpr = errors.New("not valid expression, missing '='")
var errKeyContainsSpace = errors.New("key contains space")

// Load loads env files and set environment variables.
// If not specify files, it will try to load .env file.
func Load(names ...string) error {
	if len(names) == 0 {
		names = []string{".env"}
	}
	b, err := loadFiles(names...)
	if err != nil {
		return err
	}

	m, err := Unmarshal(b)
	if err != nil {
		return err
	}
	return LoadMap(m)
}

// LoadMap loads environments from map.
func LoadMap(settings map[string]string) error {
	for k, v := range settings {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadFiles(names ...string) ([]byte, error) {
	var s bytes.Buffer
	for _, name := range names {
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}

		_, err = s.ReadFrom(f)
		if err != nil {
			return nil, err
		}

		err = f.Close()
		if err != nil {
			return nil, err
		}

		s.WriteRune('\n')
	}
	return s.Bytes(), nil
}

// Unmarshal unmarshal config from bytes.
func Unmarshal(b []byte) (map[string]string, error) {
	return Read(bytes.NewReader(b))
}

// Read config from io.Reader.
func Read(r io.Reader) (map[string]string, error) {
	envMap := make(map[string]string)
	lines, err := readLines(r)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if shouldIgnore(line) {
			continue
		}

		k, v, err := parseKV(line)
		if err != nil {
			return envMap, err
		}
		envMap[k] = v
	}
	return envMap, nil
}

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func shouldIgnore(line string) bool {
	return strings.HasSuffix(line, "#") || len(line) == 0
}

func parseKV(line string) (string, string, error) {
	// remove line comment
	line = strings.Split(line, "#")[0]

	// trim 'export '
	line = strings.TrimPrefix(line, "export ")

	equalIndex := strings.Index(line, "=")
	if equalIndex == -1 {
		return "", "", errInvalidExpr
	}

	k := strings.TrimSpace(line[:equalIndex])
	if strings.Contains(k, " ") {
		return "", "", errKeyContainsSpace
	}
	v := strings.TrimSpace(line[equalIndex+1:])
	v = strings.Trim(v, "\"")
	v = strings.Trim(v, "'")
	return k, v, nil
}
