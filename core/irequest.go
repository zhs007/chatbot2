package core

type IRequest interface {
	Start(chatbot *Chatbot, msg string) (*Message, error)

	Push(msg *Message) error

	SetCharacter(character *Character) error
}
