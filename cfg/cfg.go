package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func ReadConfig(fileName string, cfg interface{}) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("readConfig: %w", err)
	}
	reader := strings.NewReader(string(b))
	if err := json.NewDecoder(reader).Decode(&cfg); err != nil {
		return fmt.Errorf("readConfig: %w", err)
	}
	return nil
}
