package core

type User struct {
	UserID        string   `yaml:"userID" json:"userID"`               // userID
	History       []string `yaml:"history" json:"history"`             // history
	CharacterName string   `yaml:"characterName" json:"characterName"` // characterName
}

func NewUser(uid string) *User {
	return &User{
		UserID: uid,
	}
}
