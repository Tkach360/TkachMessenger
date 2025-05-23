package tcpclient

import (
    "encoding/json"
    "fmt"
    "net"

    "github.com/Tkach360/TkachMessenger/internal/core/protocol"
)

type Handler func(obj []byte)

// тип клиента TCP
// предназначен для приема сообщений и обработки их
// при помощи обработчиков
type TCPClient struct {
    conn    net.Conn
    encoder *json.Encoder
    decoder *json.Decoder

    // хеш-таблица приемников, ключ - тип объекта коммуникации, значение - соответствующий приемник
    handlers map[protocol.CommunicationObjectType]Handler
    onClose  func() // функция обработки закрытия соединения
}

func NewTCPClient(address string) (*TCPClient, error) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return nil, fmt.Errorf("connection failed: %w", err)
    }

    return &TCPClient{
        conn:     conn,
        encoder:  json.NewEncoder(conn), // кодирует сообщение в JSON и отправляет
        decoder:  json.NewDecoder(conn), // декодирует пришедшее сообщение из JSON
        handlers: make(map[protocol.CommunicationObjectType]Handler),
    }, nil
}

// создание клиента из существующего соединения
func NewTCPClientFromConn(conn net.Conn) *TCPClient {
    return &TCPClient{
        conn:     conn,
        encoder:  json.NewEncoder(conn),
        decoder:  json.NewDecoder(conn),
        handlers: make(map[protocol.CommunicationObjectType]Handler),
    }
}

func (c *TCPClient) SetOnClose(fn func()) {
    c.onClose = fn
}

// добавление нового обработчика
func (c *TCPClient) RegisterHandler(
    objtype protocol.CommunicationObjectType,
    handler Handler,
) {
    c.handlers[objtype] = handler
}

// отправка сообщения
// возвращает ошибку в случае неудачи
func (c *TCPClient) Send(msg protocol.Message) error {
    return c.encoder.Encode(msg)
}

// отправка сообщения
// возвращает ошибку в случае неудачи
func (c *TCPClient) SendAsCommunicationObject(t protocol.CommunicationObjectType, msg interface{}) error {

    content, _ := json.Marshal(msg)

    obj := protocol.CommunicationObject{
        Type:    t,
        Content: content,
    }

    return c.encoder.Encode(obj)
}

func (c *TCPClient) SendAuthRequest(userID string, password []byte) error {

    auth := protocol.AuthRequest{
        UserID:   userID,
        Password: password,
    }

    byteAuth, _ := json.Marshal(auth)

    obj := protocol.CommunicationObject{
        Type:    protocol.AUTH_REQUEST,
        Content: byteAuth,
    }

    c.encoder.Encode(obj)

    return nil
}

// установка прослушивания
func (c *TCPClient) Listen() {
    for {
        var obj protocol.CommunicationObject
        if err := c.decoder.Decode(&obj); err != nil {
            fmt.Println("Connection closed:")
            if c.onClose != nil {
                c.onClose()
            }
            return
        }

        // запускаю обработчик соответствующий типу объекта коммуникации
        if handler, ok := c.handlers[obj.Type]; ok {
            go handler(obj.Content)
        }
    }
}

// закрытие соединения
func (c *TCPClient) Close() {
    c.conn.Close()
}
