package protocol

// тип объекта коммуникации
type CommunicationObjectType int16

const (
	// запрос на аутентификацию
	AUTH_REQUEST CommunicationObjectType = iota

	// ответ на аутентификацию
	AUTH_RESPONSE

	// пользовательское сообщение
	MESSAGE
)
