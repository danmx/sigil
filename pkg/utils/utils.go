package utils

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// StringMapStringType implements cli.GenericType
type StringMapStringType struct {
	Map map[string]string
}

// Set Map value
func (m *StringMapStringType) Set(value string) error {
	tagsMap, err := stringTagsToMap(value)
	if err != nil {
		return err
	}
	m.Map = tagsMap
	return nil
}

// String return string representation of Map value
func (m *StringMapStringType) String() string {
	list := make([]string, 0, len(m.Map))
	for key, value := range m.Map {
		list = append(list, key+"="+value)
	}
	return strings.Join(list, ",")
}

func stringTagsToMap(value string) (map[string]string, error) {
	tagsMap := make(map[string]string)
	keyValuePairs := strings.Split(value, ",")
	for _, pair := range keyValuePairs {
		splittedPair := strings.Split(pair, "=")
		if len(splittedPair) != 2 {
			log.WithFields(log.Fields{
				"keyValuePairs": keyValuePairs,
				"pair":          pair,
				"splittedPair":  splittedPair,
			}).Error("wrong format of a key-value pair")
			return nil, fmt.Errorf("wrong format of a key-value pair: %s", pair)
		}
		tagsMap[splittedPair[0]] = splittedPair[1]
	}
	log.WithFields(log.Fields{
		"Tags": tagsMap,
	}).Debug("Tags Map")
	return tagsMap, nil
}
