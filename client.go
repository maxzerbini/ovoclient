package ovoclient

import (
	"encoding/json"
	"errors"
	"github.com/maxzerbini/ovoclient/model"
	"strconv"
	"sync"
)

const (
	maxServer = 128
)

// OVO Client can connect to a OVO cluster and operates with OVO server APIs.
// OVO Client is thread safe and can be shared from gorutines.
type Client struct {
	topology    *model.OvoTopology
	clients     map[string]*Session
	clientsHash map[int32]*Session
	config      *Configuration
	mux         *sync.RWMutex
}

// Create a client loading the configuration from the default path.
func NewClient() *Client {
	client := &Client{clients: make(map[string]*Session, 0), clientsHash: make(map[int32]*Session, 128), mux: new(sync.RWMutex)}
	// load configuration from default path
	client.init()
	return client
}

// Create a client using the input configuration.
func NewClientFromConfig(config *Configuration) *Client {
	client := &Client{clients: make(map[string]*Session, 0), clientsHash: make(map[int32]*Session, 128), mux: new(sync.RWMutex)}
	client.config = config
	client.init()
	return client
}

// Create a client reading the configuration file from the config-path.
func NewClientFromConfigPath(configpath string) *Client {
	client := &Client{clients: make(map[string]*Session, 0), clientsHash: make(map[int32]*Session, 128), mux: new(sync.RWMutex)}
	// load configuration
	client.config = LoadConfiguration(configpath)
	client.init()
	return client
}

// init the client
func (c *Client) init() {
	// get topology
	for _, node := range c.config.ClusterNodes {
		s := &Session{}
		res := model.OvoResponseTopology{}
		resp, err := s.Get(createTopologyEndpoint(node.Host, node.Port), nil, &res, nil)
		if err != nil {
			logInfof("Connection to %s:%s failed due to %v.\r\n", node.Host, node.Port, err)
		} else {
			if resp.Status() == 200 {
				c.topology = &res.Data
				logInfof("Connection to %s:%s done: reading topology...\r\n", node.Host, node.Port)
				break
			}
		}
	}
	c.rebuildClients()
}

// Rebuild the client map.
func (c *Client) rebuildClients() {
	c.mux.Lock()
	defer c.mux.Unlock()
	// create inner clients
	for _, node := range c.topology.Nodes {
		s := &Session{}
		s.SetNode(node)
		for _, hash := range node.HashRange {
			c.clientsHash[int32(hash)] = s
		}
		c.clients[node.Name] = s
	}
}

// Check cluster topology.
func (c *Client) checkTopology(topology model.OvoTopology) {
	// get topology
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, node := range topology.Nodes {
		s := &Session{}
		res := model.OvoResponseTopology{}
		resp, err := s.Get(createTopologyEndpoint(node.Host, strconv.Itoa(node.Port)), nil, &res, nil)
		if err != nil {
			logInfof("Connection to %s:%d failed due to %v.\r\n", node.Host, node.Port, err)
		} else {
			if resp.Status() == 200 {
				c.topology = &res.Data
				logInfof("Connection to %s:%d done: reading topology...\r\n", node.Host, node.Port)
				break
			}
		}
	}
}

// Check cluster topology and rebuild the client map.
func (c *Client) checkCluster() {
	c.checkTopology(*c.topology)
	c.rebuildClients()
}

// Get session from client map.
func (c *Client) getSessionFromHash(hash int32) *Session {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.clientsHash[hash]
}

// Put data in raw form into the OVO storage.
func (c *Client) PutRawData(key string, data []byte, ttl int) (*model.OvoResponse, error) {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	mdata := &model.OvoKVRequest{Key: key, Data: data, Hash: hash, TTL: ttl}
	resp := &model.OvoResponse{}
	if s != nil {
		_, err := s.Post(createKeyStorageEndpoint(s.node.Host, s.port), mdata, resp, nil)
		if err != nil {
			// try post on twins
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					_, errt := st.Post(createKeyStorageEndpoint(st.node.Host, st.port), mdata, resp, nil)
					if errt == nil {
						c.checkCluster()
						return resp, nil
					}
				}
			}
			c.checkCluster()
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) Put(key string, data interface{}, ttl int) (*model.OvoResponse, error) {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	bdata, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mdata := &model.OvoKVRequest{Key: key, Data: bdata, Hash: hash, TTL: ttl}
	resp := &model.OvoResponse{}
	if s != nil {
		_, err := s.Post(createKeyStorageEndpoint(s.node.Host, s.port), mdata, resp, nil)
		if err != nil {
			// try post on twins
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					_, errt := st.Post(createKeyStorageEndpoint(st.node.Host, st.port), mdata, resp, nil)
					if errt == nil {
						c.checkCluster()
						return resp, nil
					}
				}
			}
			c.checkCluster()
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) GetRawData(key string) ([]byte, error) {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	resp := &model.OvoResponse{Data: &model.OvoKVResponse{}}
	if s != nil {
		rs, err := s.Get(createGetKeyStorageEndpoint(s.node.Host, s.port, key), nil, resp, nil)
		if err != nil {
			// try get data from twins
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					rs, err := st.Get(createGetKeyStorageEndpoint(st.node.Host, st.port, key), nil, resp, nil)
					if err == nil {
						if resp != nil && rs.status == 200 {
							return resp.Data.(*model.OvoKVResponse).Data, nil
						}
					}
				}
			}
			c.checkCluster()
			return nil, errors.New("Key not found.")
		}
		if resp != nil && rs.status == 200 {
			return resp.Data.(*model.OvoKVResponse).Data, nil
		} else if rs.status == 404 {
			return nil, errors.New("Key not found.")
		}
		return nil, errors.New("Invalid data.")
	}
	return nil, errors.New("Node not found.")
}

func (c *Client) Get(key string, data interface{}) (*model.OvoResponse, error) {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	resp := &model.OvoResponse{Data: &model.OvoKVResponse{}}
	if s != nil {
		rs, err := s.Get(createGetKeyStorageEndpoint(s.node.Host, s.port, key), nil, resp, nil)
		if err != nil {
			// try get data from twins
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					rs, err := st.Get(createGetKeyStorageEndpoint(st.node.Host, st.port, key), nil, resp, nil)
					if err == nil {
						if rs.status == 200 {
							err = json.Unmarshal(resp.Data.(*model.OvoKVResponse).Data, data)
							return resp, err
						}
					}
				}
			}
			c.checkCluster()
			return nil, errors.New("Key not found.")
		}
		if rs.status == 200 {
			err = json.Unmarshal(resp.Data.(*model.OvoKVResponse).Data, data)
			return resp, err
		} else if rs.status == 404 {
			return nil, errors.New("Key not found.")
		}
		return nil, errors.New("Invalid data.")
	}
	return nil, errors.New("Node not found.")
}

// Give the number of object store in every node (also replicated object are counted) and the "TotalCount".
func (c *Client) Count() map[string]int64 {
	var count int64
	counters := make(map[string]int64, len(c.topology.Nodes))
	for _, node := range c.topology.Nodes {
		resp := &model.OvoResponse{Data: new(int64)}
		s := c.clients[node.Name]
		rs, err := s.Get(createKeyStorageEndpoint(s.node.Host, s.port), nil, resp, nil)
		if err == nil {
			if rs.status == 200 {
				counters[node.Name] = *resp.Data.(*int64)
				count += counters[node.Name]
			}
		}
	}
	counters["TotalCount"] = count
	return counters
}

// Get the list of all the keys.
func (c *Client) Keys() []string {
	keys := make(map[string]bool)
	for _, node := range c.topology.Nodes {
		resp := &model.OvoResponse{Data: &model.OvoKVKeys{}}
		s := c.clients[node.Name]
		rs, err := s.Get(createKeysEndpoint(s.node.Host, s.port), nil, resp, nil)
		if err == nil {
			if rs.status == 200 {
				for _, k := range resp.Data.(*model.OvoKVKeys).Keys {
					keys[k] = true
				}
			}
		}
	}
	klist := make([]string, 0, 0)
	for k, _ := range keys {
		klist = append(klist, k)
	}
	return klist
}

// Put data in raw form into the OVO storage.
func (c *Client) Delete(key string) (*model.OvoResponse, error) {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	resp := &model.OvoResponse{}
	if s != nil {
		_, err := s.Delete(createGetKeyStorageEndpoint(s.node.Host, s.port, key), nil, resp, nil)
		if err != nil {
			// try delete data calling twins
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					_, errt := st.Delete(createGetKeyStorageEndpoint(st.node.Host, st.port, key), nil, resp, nil)
					if errt == nil {
						c.checkCluster()
						return resp, nil
					}
				}
			}
			c.checkCluster()
			return nil, err
		}
	}
	return resp, nil
}
