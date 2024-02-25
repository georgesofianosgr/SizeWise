package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/georgesofianosgr/sizewise/pkg/storage"
)

var ErrorInvalidConfig = errors.New("Config file is invalid")

var FileName = "sizewiserc.json"

type Config struct {
	Version  float64
	Storages []storage.Storage
}

type configRaw struct {
	Version  float64           `json:"version"`
	Storages []json.RawMessage `json:"storages"`
}

func NewFromFle(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return &Config{}, err
	}
	defer file.Close()

	conf := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode((&conf))
	if err != nil {
		return &Config{}, err
	}
	return &conf, nil
}

func (c *Config) UnmarshalJSON(data []byte) error {
	cRaw := configRaw{}
	err := json.Unmarshal(data, &cRaw)
	if err != nil {
		return err
	}

	c.Version = cRaw.Version

	for _, sraw := range cRaw.Storages {
		storageBase := storage.Base{
			Cache: true,
		}
		err := json.Unmarshal(sraw, &storageBase)
		if err != nil {
			return err
		}

		switch storageBase.Type {
		case "local":
			localStorage := storage.Local{}
			err := json.Unmarshal(sraw, &localStorage)
			if err != nil {
				return err
			}
			c.Storages = append(c.Storages, localStorage)
		case "s3":
			s3Storage := storage.S3{}
			err := json.Unmarshal(sraw, &s3Storage)
			if err != nil {
				return err
			}
			c.Storages = append(c.Storages, s3Storage)
		default:
			return fmt.Errorf("%w, unknown storage type: %s", ErrorInvalidConfig, storageBase.Type)
		}
	}
	return nil
}
