package decoder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dodo/dodo-build/pkg/types"
)

func (d *decoder) DecodeSSHAgents(name string, config interface{}) ([]*types.SshAgent, error) {
	result := []*types.SshAgent{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.Bool:
		decoded, err := d.DecodeBool(name, config)
		if err != nil {
			return result, err
		}
		if decoded {
			result = append(result, &types.SshAgent{})
		}
	case reflect.String, reflect.Map:
		decoded, err := d.DecodeSSHAgent(name, config)
		if err != nil {
			return result, err
		}
		result = append(result, decoded)
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := d.DecodeSSHAgent(name, v)
			if err != nil {
				return result, err
			}
			result = append(result, decoded)
		}
	}
	return result, nil
}

func (d *decoder) DecodeSSHAgent(name string, config interface{}) (*types.SshAgent, error) {
	var result types.SshAgent
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeString(name, t.String())
		if err != nil {
			return nil, err
		}
		switch values := strings.SplitN(decoded, "=", 2); len(values) {
		case 0:
			return nil, fmt.Errorf("empty assignment in '%s'", name)
		case 1:
			return nil, fmt.Errorf("empty identity file in '%s'", name)
		case 2:
			return &types.SshAgent{Id: values[0], IdentityFile: values[1]}, nil
		default:
			return nil, fmt.Errorf("too many values in '%s'", name)
		}
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "id":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Id = decoded
			case "file":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.IdentityFile = decoded
			}
		}
	}
	return &result, nil
}
