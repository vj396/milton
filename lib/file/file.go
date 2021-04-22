package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func checkIfFileExists(filename string) error {
	f, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("required file: %q does not exist, err: %+v", filename, err)
	}
	if f.IsDir() {
		return fmt.Errorf("required file: %q is a directory", filename)
	}
	return nil
}

func ReadYAMLFile(filename string, data interface{}) error {
	if err := checkIfFileExists(filename); err != nil {
		return err
	}
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.UnmarshalStrict(f, data)
	if err != nil {
		return fmt.Errorf("cannot unmarshal data, err: %+v", err)
	}
	return nil
}

func ReadJSONFile(filename string, data interface{}) error {
	if err := checkIfFileExists(filename); err != nil {
		return err
	}
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(f, data)
	if err != nil {
		return fmt.Errorf("cannot unmarshal data, err: %+v", err)
	}
	return nil
}
