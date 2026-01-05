package database

// Database defines all of the db operation
type Database interface {
	CountAll() (int, error)                                 //Get db records count
	DeviceTokenByKey(key string) (string, error)            //Get specified device's token
	SaveDeviceTokenByKey(key, token string) (string, error) //Create or update specified devices's token
	DeleteDeviceByKey(key string) error                     //Delete specified device
	GetDeviceKeyByToken(token string) (string, error)       //Get device key by token
	SaveQiuQiuMessage(msg interface{}) error                //Save QiuQiu message to database
	GetQiuQiuMessages(token string) ([]interface{}, error)  //Get QiuQiu messages by device token
	Close() error                                           //Close the database
}
