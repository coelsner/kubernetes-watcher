package teams

import "testing"

func TestNewPodMessage(t *testing.T) {
	content := Content{
		Color: "Test", Title: "test", Text: "test",
	}

	buf, err := NewPodMessage(content)
	if err != nil {
		t.Errorf("Should not return an error: %v", err)
	}

	if buf.Len() == 0 {
		t.Errorf("Result should not be empty")
	}
}
