package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"unicode/utf8"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TG_API_TOKEN string
	Tasks        []Task
}

type TaskParamType string

const (
	ParamTypeText   TaskParamType = "text"
	ParamTypeNumber TaskParamType = "number"
	ParamTypeIP     TaskParamType = "ip"
)

var (
	ErrUnknownParamType = errors.New("unknown param type")
	ErrParamTooLong     = errors.New("param too long")
	ErrNonNumber        = errors.New("non number")
	ErrNonText          = errors.New("non text")
	ErrNonIP            = errors.New("non ip")
)

type Task struct {
	Name    string  `yaml:"name"`
	Command string  `yaml:"command"`
	Params  []Param `yaml:"params"`
}

type Param struct {
	Name string        `yaml:"name" validate:"required|max_len:20"`
	Type TaskParamType `yaml:"type" validate:"required|enum"`
}

type ParamValue struct {
	param *Param
	Value string
}

func (t *Task) ParseParamValues(msg string) ([]ParamValue, error) {
	values := make([]ParamValue, 0)

	ptrn, err := regexp.Compile(`(?:^| )(\"(?:[^\"]+|\"\")*\"|[^ ]*)`)
	if err != nil {
		return nil, err
	}

	for i, value := range ptrn.FindAllStringSubmatch(msg, -1) {
		if err := t.Params[i].Validate(value[1]); err != nil {
			return nil, fmt.Errorf("param %s: %w", t.Params[i].Name, err)
		}
		values = append(values, ParamValue{
			param: &t.Params[i],
			Value: value[1],
		})
	}
	return values, nil
}

func (tp *Param) Validate(v string) error {
	switch tp.Type {
	case ParamTypeText:
		if utf8.RuneCountInString(v) > 64 {
			return ErrParamTooLong
		}

	case ParamTypeNumber:
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			return ErrNonNumber
		}

	case ParamTypeIP:
		ip := net.ParseIP(v)
		_, _, err := net.ParseCIDR(v)

		if ip.To4() == nil && err != nil {
			return ErrNonIP
		}

	default:
		return ErrUnknownParamType
	}

	return nil
}

func Get() *Config {
	token := os.Getenv("TELEGRAM_APITOKEN")
	if token == "" {
		return nil
	}

	file, err := os.Open("config.yaml")
	if err != nil {
		return nil
	}

	defer file.Close()

	var tasks []Task
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return nil
	}

	return &Config{
		TG_API_TOKEN: token,
		Tasks: tasks,
	}
}
