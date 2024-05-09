package core

type BasicRequest struct {
	History   []*Message `yaml:"history" json:"history"`     // history
	Character *Character `yaml:"character" json:"character"` // character
}

func (req *BasicRequest) SetCharacter(character *Character) error {
	req.Character = character

	return nil
}

func (req *BasicRequest) Push(msg *Message) error {
	req.History = append(req.History, msg)

	return nil
}
