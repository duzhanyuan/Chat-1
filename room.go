package main

import (
	"container/list"
	"sync"
	"sort"
	"fmt"
)

//clientList is a mutex enhanced linked list of clients.
type clientList struct {
	*list.List
	*sync.Mutex
}

//NewClientList returns a pointer to an empty clientList.
func NewClientList () *clientList {
	return &clientList{list.New(), new(sync.Mutex)}
}

//Add adds the object c to the back of the list.
func (c *clientList) Add(cl Client) {
	c.Lock()
	c.PushBack(cl)
	c.Unlock()
}

//Rem removes all clients from the list that are equal to c.
func (c *clientList) Rem(cl Client) bool {
	c.Lock()
	found := false
	for i,x  := c.Front(), c.Front(); i != nil; {
		if other, ok := i.Value.(Client); ok {
			if cl.Equals(other) {
				x = i
				i = i.Next()
				c.Remove(x)
				found = true
			} else {
				i = i.Next()
			}
		} else {
			i = i.Next()
		}
	}
	c.Unlock()
	return found
}


//Who returns a []string with all the names of the clients in the list sorted.
func (c *clientList) Who() []string {
	clist := make([]string,0,0)
	for i:= c.Front();i != nil;i = i.Next() {
		clist = append(clist, i.Value.(Client).Name())
	}
	sort.Strings(clist)
	return clist
}


//Room is a room name and a linked list of clients in the room.
type Room struct {
	name string
	clients *clientList
	messages *messageList
}

//NewRoom creates a room with name.
func NewRoom (name string) *Room {
	newRoom := new(Room)
	newRoom.name = name
	newRoom.clients = NewClientList()
	newRoom.messages = newMessageList()
	return newRoom
}

//Equals returns true if the rooms have the same name.
func (rm *Room) Equals (other Client) bool{
	if c, ok := other.(*Room);ok {
		return rm.Name() == c.Name()
	}
	return false
}

//Name returns the name of the room.
func (rm *Room) Name() string {
	return rm.name
}

//Who returns a slice of the names of all the clients in the rooms client list.
func (rm *Room) Who () []string {
	return rm.clients.Who()
}

//Remove removes a client from the room.
func (rm *Room) Remove (cl Client) bool {
	return rm.clients.Rem(cl)
}

//Add adds a client to a room.
func (rm *Room) Add (cl Client) {
	rm.clients.Add(cl)
}

//Tell sends a string to the room from the server.
func (rm Room) Tell(s string) {
	msg := serverMessage{s}
	rm.Send(msg)
}

//Send puts the message into each client in the room's recieve function.
func (rm *Room) Send (m Message) {
	rm.clients.Lock()
	for i := rm.clients.Front(); i != nil; i = i.Next() {
		i.Value.(Client).Recieve(m)
	}
	rm.clients.Unlock()
	rm.messages.Lock()
	rm.messages.PushBack(m)
	rm.messages.Unlock()
}

//Recieve passes messages the room recieves to all clients in the room's client list.
func (rm *Room) Recieve (m Message) {
	rm.clients.Lock()
	for i := rm.clients.Front(); i != nil; i = i.Next() {
		i.Value.(Client).Recieve(m)
	}
	rm.clients.Unlock()
	rm.messages.Lock()
	rm.messages.PushBack(m)
	rm.messages.Unlock()
}


//IsEmpty returns true if the room is empty.
func (rm *Room) IsEmpty() bool {
	if rm.clients.Front() == nil {
		return true
	}
	return false
}

//GetMessages gets the messages from the room message list and returns them as a []string.
func (rm Room) GetMessages() []string {
	m := make([]string,rm.messages.Len(), rm.messages.Len())
	for i,x := rm.messages.Front(), 0; i != nil; i,x = i.Next(),x+1 {
		m[x] = fmt.Sprint(i.Value)
	}
	return m
}
