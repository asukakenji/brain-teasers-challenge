package lib

import (
	"container/list"
	"fmt"
	"sync/atomic"
	"time"
)

type MessageHash struct {
	MessageId uint64
	Message   string
	isDeleted bool
}

type Queue struct {
	lastMessageId *uint64
	idToHashMap   map[uint64]*MessageHash
	idList        *list.List
}

func NewQueue() *Queue {
	q := new(Queue)
	q.lastMessageId = new(uint64)
	q.idToHashMap = make(map[uint64]*MessageHash)
	q.idList = list.New()
	return q
}

func (q *Queue) Add(message string) (id uint64) {
	messageId := atomic.AddUint64(q.lastMessageId, 1)
	q.idToHashMap[messageId] = &MessageHash{messageId, message, false}
	q.idList.PushBack(messageId)
	return 0
}

func (q *Queue) View() *MessageHash {
	messageHash := q.idList.Remove(q.idList.Front()).(*MessageHash)
	go func() {
		select {
		case <-time.After(time.Second):
			if !messageHash.isDeleted {
				q.idList.PushFront(messageHash.MessageId)
			}
		}
	}()
	return messageHash
}

func (q *Queue) Remove(id uint64) bool {
	return false
}

func (q *Queue) PrintQueue(index int) {
	i := 0
	messageHash := q.View()
	output := ""
	for messageHash != nil {
		output += messageHash.Message
		output += " "
		if i == index {
			q.Remove(messageHash.MessageId)
		}
		i++
		messageHash = q.View()
	}
	if output != "" {
		fmt.Println(output)
	}
}
