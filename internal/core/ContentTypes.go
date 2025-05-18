package core

// тип контента в сообщении
type ContentType int16

// при добавлении новых типов контента придется добавлять сюда их идентификаторы
const (
	TEXT_CONTENT  ContentType = iota
	IMAGE_CONTENT             // возможно реализую
)
