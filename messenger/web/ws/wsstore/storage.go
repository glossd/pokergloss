package wsstore

import (
	"fmt"
	"sync"
)

var ErrNotOnline = fmt.Errorf("user is not online")

var userStore = &sync.Map{}

func GetClient(userID string) (*Client, error) {
	data, ok := userStore.Load(userID)
	if !ok {
		return nil, ErrNotOnline
	}
	client := data.(*Client)
	return client, nil
}

func addUser(client *Client) {
	userStore.Store(client.Identity.UserId, client)
}

func RemoveUser(client *Client) {
	userStore.Delete(client.Identity.UserId)
	client.Close()
}
