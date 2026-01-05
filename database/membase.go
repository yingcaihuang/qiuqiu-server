package database

import (
	"fmt"
	"strings"
)

var (
	cacheKey         = "MemoryBaseKey"
	cacheDeviceToken = ""
)

type MemBase struct {
}

func NewMemBase() Database {
	return &MemBase{}
}

func (d *MemBase) CountAll() (int, error) {
	return 1, nil
}

func (d *MemBase) DeviceTokenByKey(key string) (string, error) {
	if cacheKey == key && cacheDeviceToken != "" {
		return cacheDeviceToken, nil
	}
	return "nil", fmt.Errorf("key not found")
}

func (d *MemBase) SaveDeviceTokenByKey(key, token string) (string, error) {
	if key != "" && key != cacheKey {
		return "", fmt.Errorf("key not found")
	}
	// Deep copy prevents Fiber memory overwrite bugs.
	cacheDeviceToken = strings.Clone(token)
	return key, nil
}

func (d *MemBase) DeleteDeviceByKey(key string) error {
	if key != "" && key != cacheKey {
		return fmt.Errorf("key not found")
	}
	cacheDeviceToken = ""
	return nil
}

func (d *MemBase) GetDeviceKeyByToken(token string) (string, error) {
	if cacheDeviceToken == token {
		return cacheKey, nil
	}
	return "", fmt.Errorf("token not found")
}

func (d *MemBase) SaveQiuQiuMessage(msg interface{}) error {
	// In-memory implementation doesn't persist messages
	return nil
}

func (d *MemBase) GetQiuQiuMessages(token string) ([]interface{}, error) {
	// Return empty message list for in-memory implementation
	return []interface{}{}, nil
}

func (d *MemBase) Close() error {
	return nil
}
