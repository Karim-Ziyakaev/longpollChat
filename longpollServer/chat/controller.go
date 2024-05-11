package chat

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Controller struct {
	messages map[string]chan Message
	mutex    sync.RWMutex
	db       *gorm.DB
}

func NewController(db *gorm.DB) *Controller {
	err := db.AutoMigrate(&Message{})
	if err != nil {
		log.Fatal(err)
	}
	return &Controller{
		messages: make(map[string]chan Message),
		db:       db,
	}
}

func (c *Controller) Send(msg Message) error {
	c.mutex.RLock()
	ch, ok := c.messages[msg.To.String()]
	c.mutex.RUnlock()
	if !ok {
		c.mutex.Lock()
		c.messages[msg.To.String()] = make(chan Message, 50)
		ch = c.messages[msg.To.String()]
		c.mutex.Unlock()
	}

	ch <- msg
	err := c.db.Create(&msg).Error
	return err
}

func (c *Controller) GetMessages(userID string) ([]Message, error) {
	var messages []Message

	err := c.db.Where("\"to\" = ?", userID).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Controller) GetMessage(userID string) (Message, error) {
	var msg Message

	c.mutex.RLock()
	ch, ok := c.messages[userID]
	c.mutex.RUnlock()

	if !ok {
		c.mutex.Lock()
		c.messages[userID] = make(chan Message, 50)
		ch = c.messages[userID]
		c.mutex.Unlock()
	}

	select {
	case msg = <-ch:
		return msg, nil
	case <-time.After(1 * time.Minute):
		return msg, errors.New("timeout")
	}
}
