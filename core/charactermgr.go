package core

import (
	"os"

	"github.com/zhs007/goutils"
	"gopkg.in/yaml.v3"
)

type CharacterMgr struct {
	MapCharacters    map[string]*Character `yaml:"characters" json:"characters"` // characters
	DefaultCharacter string                `yaml:"default" json:"default"`       // default
}

func (mgr *CharacterMgr) Get(character string) *Character {
	c, isok := mgr.MapCharacters[character]
	if isok {
		return c
	}

	return nil
}

func (mgr *CharacterMgr) GetDefault() *Character {
	return mgr.MapCharacters[mgr.DefaultCharacter]
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

	for k, v := range mgr.MapCharacters {
		v.Name = k
	}

	return mgr, nil
}
