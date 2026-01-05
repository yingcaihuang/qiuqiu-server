package database

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

// QiuQiuMessageRecord stores a message with metadata
type QiuQiuMessageRecord struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	URL       string    `json:"url,omitempty"`
	Timestamp int64     `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	qiuqiuMessagesBucket = "qiuqiu_messages"
	deviceTokenBucket    = "device_tokens_reverse"
)

// GetDeviceKeyByToken retrieves device key using device token
// This implements reverse lookup from token to device key
func (d *BboltDB) GetDeviceKeyByToken(token string) (string, error) {
	var key string
	err := db.View(func(tx *bbolt.Tx) error {
		reverseBucket := tx.Bucket([]byte(deviceTokenBucket))
		if reverseBucket == nil {
			return fmt.Errorf("device token reverse lookup bucket not found")
		}

		if bs := reverseBucket.Get([]byte(token)); bs == nil {
			return fmt.Errorf("failed to get device key for token [%s]", token)
		} else {
			key = string(bs)
			return nil
		}
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

// SaveQiuQiuMessage saves a QiuQiu message to the database
func (d *BboltDB) SaveQiuQiuMessage(msgInterface interface{}) error {
	// Type assertion to handle the message
	var token, title, message, url string
	var timestamp int64

	switch msg := msgInterface.(type) {
	case map[string]interface{}:
		token, _ = msg["token"].(string)
		title, _ = msg["title"].(string)
		message, _ = msg["message"].(string)
		url, _ = msg["url"].(string)
		if t, ok := msg["timestamp"].(int64); ok {
			timestamp = t
		} else if t, ok := msg["timestamp"].(float64); ok {
			timestamp = int64(t)
		}
	default:
		// Try JSON marshaling as fallback
		if data, err := json.Marshal(msgInterface); err != nil {
			return fmt.Errorf("failed to marshal message: %v", err)
		} else {
			var m map[string]interface{}
			if err := json.Unmarshal(data, &m); err != nil {
				return fmt.Errorf("failed to unmarshal message: %v", err)
			}
			token, _ = m["token"].(string)
			title, _ = m["title"].(string)
			message, _ = m["message"].(string)
			url, _ = m["url"].(string)
			if t, ok := m["timestamp"].(int64); ok {
				timestamp = t
			} else if t, ok := m["timestamp"].(float64); ok {
				timestamp = int64(t)
			}
		}
	}

	if token == "" {
		return fmt.Errorf("message token is empty")
	}

	if message == "" {
		return fmt.Errorf("message content is empty")
	}

	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	record := QiuQiuMessageRecord{
		ID:        fmt.Sprintf("%s-%d", token, timestamp),
		Token:     token,
		Title:     title,
		Message:   message,
		URL:       url,
		Timestamp: timestamp,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal message record: %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		// Create messages bucket if not exists
		bucket, err := tx.CreateBucketIfNotExists([]byte(qiuqiuMessagesBucket))
		if err != nil {
			return err
		}

		// Save message by token-timestamp key
		return bucket.Put([]byte(record.ID), data)
	})

	return err
}

// GetQiuQiuMessages retrieves all messages for a given device token
func (d *BboltDB) GetQiuQiuMessages(token string) ([]interface{}, error) {
	var messages []interface{}

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(qiuqiuMessagesBucket))
		if bucket == nil {
			// Bucket doesn't exist yet, return empty list
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var record QiuQiuMessageRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return nil // Skip malformed records
			}

			// Filter by token
			if record.Token == token {
				messages = append(messages, record)
			}

			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return messages, nil
}
