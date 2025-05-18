package protocol

import "encoding/json"

// тип данных для сообщения, которое передается по сети
// внутри может содержать как служебное сообщение для приложения
// так и сообщение с контентом
type CommunicationObject struct {
	Type int16 `json:"Type"`

	// json.RawMessage позволяет отложить десериализацию, то есть при
	// десериализации ProtocolMessage содержимое Content не будет десериализовано,
	// производить десериализацию буду после определения типа Type
	Content json.RawMessage `json:"Content"`
}
