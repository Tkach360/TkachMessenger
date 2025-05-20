package tcpclient

import (
    "encoding/json"
    "fmt"
    "net"

    "github.com/Tkach360/TkachMessenger/internal/core/protocol"
)

// структура функции-приемника сообщений
type MessageHandler func(msg protocol.Message)

// тип клиента TCP
// предназначен для приема сообщений и обработки их
// при помощи обработчиков
type TCPClient struct {
    conn     net.Conn
    encoder  *json.Encoder
    decoder  *json.Decoder
    handlers []MessageHandler
}

func NewTCPClient(address string) (*TCPClient, error) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return nil, fmt.Errorf("connection failed: %w", err)
    }

    return &TCPClient{
        conn:    conn,
        encoder: json.NewEncoder(conn), // кодирует сообщение в JSON и отправляет
        decoder: json.NewDecoder(conn), // декодирует пришедшее сообщение из JSON
    }, nil
}

// добавление нового обработчика
func (c *TCPClient) RegisterHandler(handler MessageHandler) {
    c.handlers = append(c.handlers, handler)
}

// отправка сообщения
// возвращает ошибку в случае неудачи
func (c *TCPClient) Send(msg protocol.Message) error {
    return c.encoder.Encode(msg)
}

// установка прослушивания
func (c *TCPClient) Listen() {
    for {
        var msg protocol.Message
        if err := c.decoder.Decode(&msg); err != nil {
            fmt.Println("Connection closed:", err)
            return
        }

        for _, handler := range c.handlers {
            go handler(msg)
        }
    }
}

// закрытие соединения
func (c *TCPClient) Close() {
    c.conn.Close()
}
