package ovoclient

import (
	"strconv"
	"log"
)

import(
	"github.com/maxzerbini/ovoclient/model"
	"sync"
)

type Client struct {
	topology *model.OvoTopology
	clients map[string]*Session
	clientsHash map[int]*Session
	config *Configuration
	mux *sync.RWMutex
}

func NewClient() *Client {
	client := &Client{clients:make(map[string]*Session,0),clientsHash:make(map[int]*Session,128), mux:new(sync.RWMutex)}
	// load configuration from default path
	client.init()
	return client
}

func NewClientFromConfig(config *Configuration) *Client {
	client := &Client{clients:make(map[string]*Session,0),clientsHash:make(map[int]*Session,128), mux:new(sync.RWMutex)}
	client.config = config
	client.init()
	return client
}

func NewClientFromConfigPath(configpath string) *Client {
	client := &Client{clients:make(map[string]*Session,0),clientsHash:make(map[int]*Session,128), mux:new(sync.RWMutex)}
	// load configuration
	client.config = LoadConfiguration(configpath)
	client.init()
	return client
}

func (c *Client) init(){
	// get topology
	for _, node := range c.config.ClusterNodes {
		s := &Session{}
		res := model.OvoResponseTopology{}
		resp, err := s.Get(createTopologyEndpoint(node.Host, node.Port), nil, &res, nil)
		if err != nil {
			log.Printf("Connection to %s:%s failed due to %v.\r\n", node.Host, node.Port, err)
		} else {
			if resp.Status() == 200 {
				c.topology = &res.Data
				log.Printf("Connection to %s:%s done: reading topology...\r\n", node.Host, node.Port)
				break
			}
		}	
	}
	c.rebuildClients()
}

func (c *Client) rebuildClients(){
	c.mux.Lock()
	defer c.mux.Unlock()
	// create inner clients
	for _, node := range c.topology.Nodes {
		s := &Session{}
		res := model.OvoResponseTopologyNode{}
		resp, err := s.Get(createTopologyNodeEndpoint(node.Host, strconv.Itoa(node.Port)), nil, &res, nil)
		if err != nil {
			log.Printf("Connection to %s:%s failed.\r\n", node.Host, node.Port)
		} else {
			if resp.Status() == 200 {
				s.Node = node
				log.Printf("Found node %v in topology\r\n", res.Data)
				c.clients[res.Data.Name] = s
				for _,hash := range res.Data.HashRange {
					c.clientsHash[hash] = s
				}
			}
		}	
	}
}

// Add a session
func (c *Client) addSession(){
	
}
// Remove and close session
func (c *Client) removeSession(){
	
}
// Get a Session
func (c *Client) getSession(hash int) *Session{
	return nil
}


