package main

import (
	"fmt"
	"time"

	"github.com/finb/bark-server/v2/apns"
	"github.com/gofiber/fiber/v2"
)

// QiuQiuMessage represents a message structure for QiuQiu app
type QiuQiuMessage struct {
	Token     string `json:"token" form:"token" query:"token" validate:"required"`
	Title     string `json:"title" form:"title" query:"title"`
	Message   string `json:"message" form:"message" query:"message" validate:"required"`
	URL       string `json:"url" form:"url" query:"url"`
	Timestamp int64  `json:"timestamp" form:"timestamp" query:"timestamp"`
}

func init() {
	registerRoute("qiuqiu", func(router fiber.Router) {
		// QiuQiu webhook endpoint for pushing messages
		router.Post("/qiuqiu/push", routeQiuQiuPush)
		router.Post("/api/push", routeQiuQiuPush) // Alias for easy webhook integration
		router.Get("/qiuqiu/messages/:token", routeQiuQiuGetMessages)
	})
}

// routeQiuQiuPush handles incoming push requests from webhook integrations
// Expects JSON payload with token, title, message, url, and timestamp
func routeQiuQiuPush(c *fiber.Ctx) error {
	var msg QiuQiuMessage
	if err := c.BodyParser(&msg); err != nil {
		return c.Status(400).JSON(failed(400, "invalid request body: %v", err))
	}

	// Validate required fields
	if msg.Token == "" {
		return c.Status(400).JSON(failed(400, "token is required"))
	}

	if msg.Message == "" {
		return c.Status(400).JSON(failed(400, "message is required"))
	}

	// Set default title if not provided
	if msg.Title == "" {
		msg.Title = "Alert"
	}

	// Set timestamp to current time if not provided
	if msg.Timestamp == 0 {
		msg.Timestamp = time.Now().Unix()
	}

	// Save message to database
	if err := db.SaveQiuQiuMessage(msg); err != nil {
		return c.Status(500).JSON(failed(500, "failed to save message: %v", err))
	}

	// Push notification via APNs using the device token
	code, err := pushQiuQiuNotification(msg)
	if err != nil {
		// Log error but still return success since message is saved
		return c.Status(200).JSON(data(map[string]interface{}{
			"code":    code,
			"message": "Message saved but push notification failed: " + err.Error(),
			"token":   msg.Token,
		}))
	}

	return c.Status(200).JSON(data(map[string]interface{}{
		"code":    200,
		"message": "success",
		"token":   msg.Token,
	}))
}

// routeQiuQiuGetMessages retrieves all messages for a given device token
func routeQiuQiuGetMessages(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(400).JSON(failed(400, "token is required"))
	}

	messages, err := db.GetQiuQiuMessages(token)
	if err != nil {
		return c.Status(500).JSON(failed(500, "failed to retrieve messages: %v", err))
	}

	return c.Status(200).JSON(data(map[string]interface{}{
		"messages": messages,
		"token":    token,
	}))
}

// pushQiuQiuNotification sends a push notification via APNs
func pushQiuQiuNotification(msg QiuQiuMessage) (int, error) {
	// Look up device key from device token
	// For QiuQiu, we use token as device_token directly
	deviceKey, err := db.GetDeviceKeyByToken(msg.Token)
	if err != nil {
		// If device not found, try to register it automatically
		newKey, err := db.SaveDeviceTokenByKey("", msg.Token)
		if err != nil {
			return 500, fmt.Errorf("failed to register device: %v", err)
		}
		deviceKey = newKey
	}

	// Construct APNs message
	apnsMsg := apnsMessageFromQiuQiu(msg, msg.Token)
	apnsMsg.DeviceKey = deviceKey

	// Push via APNs
	code, err := apns.Push(&apnsMsg)
	if err != nil {
		return code, err
	}

	return 200, nil
}

// apnsMessageFromQiuQiu converts a QiuQiuMessage to an APNS PushMessage
func apnsMessageFromQiuQiu(msg QiuQiuMessage, deviceToken string) apns.PushMessage {
	apnsMsg := apns.PushMessage{
		Title:       msg.Title,
		Subtitle:    "",
		Body:        msg.Message,
		Sound:       "1107",
		DeviceToken: deviceToken,
		ExtParams:   make(map[string]interface{}),
	}

	// Add URL to extra params if provided
	if msg.URL != "" {
		apnsMsg.ExtParams["url"] = msg.URL
	}

	// Add timestamp
	apnsMsg.ExtParams["timestamp"] = msg.Timestamp

	return apnsMsg
}
