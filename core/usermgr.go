package core

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"gopkg.in/yaml.v3"
)

type UserMgr struct {
	MapUsers map[string]*User `yaml:"users" json:"users"` // users
}

func (mgr *UserMgr) GetUser(uid string, mgrCharacters *CharacterMgr) *User {
	u, isok := mgr.MapUsers[uid]
	if isok {
		return u
	}

	goutils.Debug("UserMgr.GetUser:NewUser",
		slog.String("uid", uid))

	nu := NewUser(uid)
	nu.Rebuild(mgrCharacters)

	mgr.MapUsers[uid] = nu

	return nu
}

func (mgr *UserMgr) Rebuild(mgrCharacters *CharacterMgr) {
	for _, user := range mgr.MapUsers {
		user.Rebuild(mgrCharacters)
	}
}

func NewUserMgr() *UserMgr {
	return &UserMgr{
		MapUsers: make(map[string]*User),
	}
}

func LoadUserMgr(fn string) (*UserMgr, error) {
	buf, err := os.ReadFile(fn)
	if err != nil {
		return NewUserMgr(), nil
	}

	mgr := &UserMgr{}

	err = yaml.Unmarshal(buf, mgr)
	if err != nil {
		goutils.Error("LoadUserMgr:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	return mgr, nil
}
