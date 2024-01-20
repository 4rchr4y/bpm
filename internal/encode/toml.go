package encode

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/4rchr4y/godevkit/syswrap"
	"github.com/BurntSushi/toml"
)

type tomlEncoderOSWrapper interface {
	LookupEnv(key string) (string, bool)
}

type TomlEncoder struct {
	osWrap tomlEncoderOSWrapper
}

func NewTomlEncoder() *TomlEncoder {
	return &TomlEncoder{
		osWrap: new(syswrap.OsWrapper),
	}
}

const envVarString = `\$\{[A-Za-z_][A-Za-z0-9_]*\}`

var envVarPattern = regexp.MustCompile(envVarString)

func (ts *TomlEncoder) Decode(data string, value interface{}) error {
	content, err := ts.interpolate(data)
	if err != nil {
		return err
	}

	if _, err := toml.Decode(content, value); err != nil {
		return err
	}

	return nil
}

func (ts *TomlEncoder) interpolate(data string) (string, error) {
	// preliminary check for the presence of placeholders
	if !strings.Contains(data, "${") {
		return data, nil
	}

	var missingVars []string
	result := envVarPattern.ReplaceAllStringFunc(data, func(match string) string {
		envKey := strings.Clone(match[2 : len(match)-1])
		if value, exists := ts.osWrap.LookupEnv(envKey); exists {
			return value
		}

		missingVars = append(missingVars, envKey)

		return match
	})

	// check if there are any unresolved variables
	if len(missingVars) > 0 {
		return "", fmt.Errorf("environment variables not found: %s", strings.Join(missingVars, ", "))
	}

	return result, nil
}

func (ts *TomlEncoder) Encode(writer io.Writer, value interface{}) error {
	return toml.NewEncoder(writer).Encode(value)
}
