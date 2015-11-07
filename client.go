package ovoclient

import(
	"github.com/maxzerbini/ovoclient/model"
)

type Client struct {
	topology *model.OvoTopology
	clients map[string]*Session
	config *Configuration
}

func NewClient() *Client {
	client := &Client{clients:make(map[string]*Session,0)}
	client.init()
	return client
}

func (c *Client) init(){
	
}