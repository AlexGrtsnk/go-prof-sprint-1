package internal

import "testing"

func TestNewProducer(t *testing.T) {
	_, err := NewProducer("")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}

func TestNewConsumer(t *testing.T) {
	_, err := NewConsumer("")
	if err == nil {
		t.Errorf("IntMin(2, -2) = %d; want -2", err)
	}
}
