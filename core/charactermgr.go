package core

import (
	"os"

	"github.com/zhs007/goutils"
	"gopkg.in/yaml.v3"
)

type CharacterMgr struct {
	Characters []*Character `yaml:"characters" json:"characters"` // characters
}

func LoadCharacterMgr(fn string) (*CharacterMgr, error) {
	buf, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("LoadCharacterMgr:ReadFile",
			goutils.Err(err))

		return nil, err
	}

	mgr := &CharacterMgr{}

	err = yaml.Unmarshal(buf, mgr)
	if err != nil {
		goutils.Error("LoadCharacterMgr:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	return mgr, nil
}
