package core

type UserMgr struct {
	MapUsers map[string]*User `yaml:"users" json:"users"` // users
}

func (mgr *UserMgr) GetUser(uid string) *User {
	u, isok := mgr.MapUsers[uid]
	if isok {
		return u
	}

	nu := NewUser(uid)

	mgr.MapUsers[uid] = nu

	return nu
}

func NewUserMgr() *UserMgr {
	return &UserMgr{
		MapUsers: make(map[string]*User),
	}
}
