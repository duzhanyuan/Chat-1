package main

import (
	"testing"
	"time"
	"fmt"
)

var messageTests = []struct {
	name string
	text string
	time string
	result string
}{
	{"Bob","Hello World", "1:00pm", "1:00pm [Bob]: Hello World"},
	{"Fred","hi","7:00pm","7:00pm [Fred]: hi"},
	{"dsAI!kd","839skdl9","12:00am","12:00am [dsAI!kd]: 839skdl9"},
}
type MTestClient struct {
	name string
	messages *messageList
}

func (cl *MTestClient) Equals(other Client) bool {
	if cl.Name() == other.Name() {
		return true
	}
	return false
}

func (cl *MTestClient) Recieve(m Message) {
	cl.messages.PushBack(m)
}

func (cl *MTestClient) Name() string {
	return cl.name
}


func TestRestMessageString(t *testing.T) {
	msg := new(restMessage)
	var err error
	for _, tt := range messageTests {
		msg.Name = tt.name
		msg.Text = tt.text
		msg.Time, err = time.Parse("3:04pm", tt.time)
		if err != nil {
			fmt.Println("Error Parsing Time: ", err)
		}
		if msg.String() != tt.result {
			t.Errorf("msg.String() %q,%q,%v => %q, want %q", msg.Name, msg.Text, msg.Time,msg.String(), tt.result)
		}
	}
}

func MTestClientMessageString(t *testing.T) {
	msg := new(clientMessage)
	cl := new(MTestClient)
	var err error
	for _, tt := range messageTests {
		cl.name = tt.name
		msg.Sender = cl
		msg.text = tt.text
		msg.time, err = time.Parse("3:04pm", tt.time)
		if err != nil {
			fmt.Println("Error Parsing Time: ", err)
		}
		if msg.String() != tt.result {
			t.Errorf("msg.String() %q,%q,%v => %q, want %q",msg.Sender.Name(),msg.text,msg.time,msg.String(), tt.result)
		}
	}
}
