package internal

import "testing"

func TestNewProducer(t *testing.T) {
	_, err := NewProducer("")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}

func TestNewConsumer(t *testing.T) {
	_, err := NewConsumer("")
	if err == nil {
		t.Errorf("this is err = %d", err)
	}
}
