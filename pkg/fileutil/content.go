package fileutil

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/4rchr4y/bpm/constant"
)

func ReadLinesToMap(content []byte) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		result[line] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading '%s' input: %v", constant.IgnoreFileName, err)
	}

	return result, nil
}
