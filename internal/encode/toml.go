package encode

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const (
	envVarString = `\$\{[A-Za-z_][A-Za-z0-9_]*\}`
)

var (
	envVarPattern = regexp.MustCompile(envVarString)
)

type tomlEncoderOSWrapper interface {
	LookupEnv(key string) (string, bool)
}

type TomlEncoder struct {
	osWrap tomlEncoderOSWrapper
}

// func NewTomlEncoder() *TomlEncoder {
// 	return &TomlEncoder{
// 		osWrap: new(syswrap.OsWrapper),
// 	}
// }

func (ts *TomlEncoder) Encode(value interface{}) ([]byte, error) {
	return toml.Marshal(value)
}

func (ts *TomlEncoder) Decode(data string, value interface{}) error {
	content, err := ts.interpolate(data)
	if err != nil {
		return err
	}

	if err = toml.Unmarshal([]byte(content), value); err != nil {
		var decodeErr *toml.DecodeError
		if ok := errors.As(err, &decodeErr); ok {
			return errors.New(decodeErr.String())

		}

		return err
	}

	return nil
}

func (ts *TomlEncoder) interpolate(data string) (string, error) {
	if !strings.Contains(data, "${") {
		return data, nil
	}

	var missingVars []string
	var result strings.Builder

	envVarPattern.ReplaceAllStringFunc(data, func(match string) string {
		envKey := strings.Clone(match[2 : len(match)-1])
		if value, exists := ts.osWrap.LookupEnv(envKey); exists {
			result.WriteString(value)
			return ""
		}

		missingVars = append(missingVars, envKey)
		result.WriteString(match)
		return ""
	})

	if len(missingVars) > 0 {
		return "", fmt.Errorf("environment variables not found: %s", strings.Join(missingVars, ", "))
	}

	return result.String(), nil
}
