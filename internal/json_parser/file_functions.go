package internal

import (
	"encoding/json"
	"os"
)

// Event тип, содержащий строку того, что хранится в базе данных для представления в формате json
type Event struct {
	// ID номер записи в базе данных
	ID int `json:"id"`
	// ShortURL сокращенный url
	ShortURL string `json:"short_url"`
	// LongURL изначальный отправленный пользователем url
	LongURL string `json:"longURL"`
	// Token поле токена пользователя, сделавшего запрос на скоращение
	Token string `json:"token"`
	// DelFlag флаг мягкого удаления из базы данныз
	DelFlag int `json:"delFlag"`
}

// Producer тип помощи записи в файл в формате json
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// ProduceURL тип записи о url в json
type ProduceURL struct {
	// CorrelationID id посланного нам url
	CorrelationID string `json:"correlation_id"`
	// OriginalURL изначальный адрес url для сокращения
	OriginalURL string `json:"original_url"`
}

// ProduceList массив созданных url
type ProduceList []ProduceURL

// DeleteURL - url, который необходимо удалить базы данных
type DeleteURL string

// DeleteList - массив url, которые необходимо удалить базы данных
type DeleteList []DeleteURL

// NewProducer - функция producer для записи в файл в формате json
func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteEvent - событие записи
func (p *Producer) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}

// Close - событие закрытия файла
func (p *Producer) Close() error {
	return p.file.Close()
}

// Consumer - тип для чтения записей из файла
type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

// NewConsumer - функция чтения записей из файла
func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadEvent событие чтения из файла
func (c *Consumer) ReadEvent() (*Event, error) {
	event := &Event{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

// Close - событие закрытия файла
func (c *Consumer) Close() error {
	return c.file.Close()
}
