package internal

import (
	"encoding/json"
	"os"
)

type Event struct {
	ID       int    `json:"id"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"longURL"`
	Tknm     string `json:"tknm"`
	DelFlag  int    `json:"delFlag"`
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

type ProduceURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type ProduceList []ProduceURL

/*
type DeleteURL struct {
	ShortURL string `json:"-"`
}*/
type DeleteURL string
type DeleteList []DeleteURL

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

func (p *Producer) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

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

func (c *Consumer) ReadEvent() (*Event, error) {
	event := &Event{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
