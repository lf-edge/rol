package infrastructure

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"rol/app/errors"
)

func createYamlFile(fileName string, obj interface{}) error {
	yamlData, err := yaml.Marshal(obj)
	if err != nil {
		return errors.Internal.Wrap(err, "yaml marshal failed")
	}
	err = ioutil.WriteFile(fileName, yamlData, 0664)
	if err != nil {
		return errors.Internal.Wrap(err, "write to file error")
	}
	return nil
}

//SaveYamlFile saving struct to a yaml file along the given path
//
//Return:
//	error - if an error occurs, otherwise nil
func SaveYamlFile(obj interface{}, filePath string) error {
	err := createYamlFile(filePath, obj)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to save yaml file")
	}

	return nil
}

//ReadYamlFile reads the yaml file at the given path into the specified struct
//
//Return:
//	StructType - specified struct with yaml data
//	error - if an error occurs, otherwise nil
func ReadYamlFile[StructType interface{}](filePath string) (StructType, error) {
	config := new(StructType)
	if _, err := os.Stat(filePath); err != nil {
		return *config, errors.NotFound.Wrap(err, "file not found")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return *config, errors.Internal.Wrap(err, "error while open file")
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		return *config, errors.Internal.Wrap(err, "failed to decode yaml")
	}

	return *config, nil
}
