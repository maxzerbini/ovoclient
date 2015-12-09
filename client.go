package ovoclient

import (
	"encoding/json"
	"errors"
	"github.com/maxzerbini/ovoclient/model"
	"strconv"
	"sync"
	"time"
)

const (
	maxServer = 128
	minClusterCheckPeriod = 10
)

// OVO Client can connect to a OVO cluster and operates with OVO server APIs.
// OVO Client is thread safe and can be shared from gorutines.
type Client struct {
	topology    *model.OvoTopology
	clients     map[string]*Session
	clientsHash map[int32]*Session
	config      *Configuration
	mux         *sync.RWMutex
	tickChan <-chan time.Time
	doneChan chan bool
}

// Create a client loading the configuration from the default path.
func NewClient() *Client {
	client := &Client{clients: make(map[string]*Session, 0), clientsHash: make(map[int32]*Session, 128), mux: new(sync.RWMutex)}
	client.tickChan = time.NewTicker(time.Second * minClusterCheckPeriod).C
	client.doneChan = make(chan bool)
	// load configuration from default path
	client.init()
	go client.check()
	return client
}

// Create a client using the input configuration.
func NewClientFromConfig(config *Configuration) *Client {
	client := &Client{clients: make(map[string]*Session, 0), clientsHash: make(map[int32]*Session, 128), mux: new(sync.RWMutex)}
	client.config = config
	client.tickChan = time.NewTicker(time.Second * 30).C
	client.doneChan = make(chan bool)
	client.init()
	go client.check()
	return client
}

// Create a client reading the configuration file from the config-path.
func NewClientFromConfigPath(configpath string) *Client {
	client := &Client{clients: make(map[string]*Session, 0), clientsHash: make(map[int32]*Session, 128), mux: new(sync.RWMutex)}
	client.tickChan = time.NewTicker(time.Second * 30).C
	client.doneChan = make(chan bool)
	// load configuration
	client.config = LoadConfiguration(configpath)
	client.init()
	go client.check()
	return client
}

// init the client
func (c *Client) init() {
	if c.config.ClusterCheckPeriod < minClusterCheckPeriod { c.config.ClusterCheckPeriod = minClusterCheckPeriod}
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

// Check cluster periodically.
func (c *Client) check(){
	for {
        select {
        case <- c.tickChan:
            c.checkCluster()
        case <- c.doneChan:
            return
      }
    }
}

// Close the client.
func (c *Client) Close(){
	c.doneChan <- true
	
}

// Put data in raw format into the OVO storage.
// The parameter key is the string associated to the object.
// The parameter data is the array of bytes rapresenting the object.
// The parameter ttl is the time to live of the object expressed in seconds; if it's zero the object will not be removed from the storage.
func (c *Client) PutRawData(key string, data []byte, ttl int) error {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	mdata := &model.OvoKVRequest{Key: key, Data: data, Hash: hash, TTL: ttl}
	resp := &model.OvoResponse{}
	if s != nil {
		_, err := s.Post(createKeyStorageEndpoint(s.node.Host, s.port), mdata, resp, nil)
		if err != nil {
			done := true
			// try post on twins
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					_, errt := st.Post(createKeyStorageEndpoint(st.node.Host, st.port), mdata, resp, nil)
					done = done && (errt ==nil)
				}
			}
			c.checkCluster()
			if done {
				return nil
			} else {
				return err
			}
		}
		return nil
	}
	return errors.New("Node not found.")
}

// Put the object in the storage serializing it in JSON.
// The parameter ttl is the time to live of the object expressed in seconds; if it's zero the object will not be removed from the storage.
func (c *Client) Put(key string, data interface{}, ttl int) error {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	bdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	mdata := &model.OvoKVRequest{Key: key, Data: bdata, Hash: hash, TTL: ttl}
	resp := &model.OvoResponse{}
	if s != nil {
		_, err := s.Post(createKeyStorageEndpoint(s.node.Host, s.port), mdata, resp, nil)
		if err != nil {
			done := true
			// try post on twins
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					_, errt := st.Post(createKeyStorageEndpoint(st.node.Host, st.port), mdata, resp, nil)
					done = done && (errt ==nil)
				}
			}
			c.checkCluster()
			if done {
				return nil
			} else {
				return err
			}
		}
		return nil
	}
	return errors.New("Node not found.")
}

// Get a raw format rapresentation of the object stored in the OVO cluster.
func (c *Client) GetRawData(key string) ([]byte, error) {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	resp := &model.OvoResponse{Data: &model.OvoKVResponse{}}
	if s != nil {
		rs, err := s.Get(createGetKeyStorageEndpoint(s.node.Host, s.port, key), nil, resp, nil)
		if err != nil {
			// try get data from the twins
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

// Retrieve an object previously serialized in JSON.
func (c *Client) Get(key string, data interface{}) error {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	resp := &model.OvoResponse{Data: &model.OvoKVResponse{}}
	if s != nil {
		rs, err := s.Get(createGetKeyStorageEndpoint(s.node.Host, s.port, key), nil, resp, nil)
		if err != nil {
			// try get data from the twins
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					rs, err := st.Get(createGetKeyStorageEndpoint(st.node.Host, st.port, key), nil, resp, nil)
					if err == nil {
						if rs.status == 200 {
							err = json.Unmarshal(resp.Data.(*model.OvoKVResponse).Data, data)
							return err
						}
					}
				}
			}
			c.checkCluster()
			return errors.New("Key not found.")
		}
		if rs.status == 200 {
			err = json.Unmarshal(resp.Data.(*model.OvoKVResponse).Data, data)
			return err
		} else if rs.status == 404 {
			return errors.New("Key not found.")
		}
		return errors.New("Invalid data.")
	}
	return errors.New("Node not found.")
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
func (c *Client) Delete(key string) error {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	resp := &model.OvoResponse{}
	if s != nil {
		_, err := s.Delete(createGetKeyStorageEndpoint(s.node.Host, s.port, key), nil, resp, nil)
		if err != nil {
			// delete data calling all the twins
			done := true
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					_, errt := st.Delete(createGetKeyStorageEndpoint(st.node.Host, st.port, key), nil, resp, nil)
					done = done && (errt ==nil)
				}
			}
			c.checkCluster()
			if done {
				return nil
			} else {
				return err
			}
		} else {
			return nil
		}
	}
	return errors.New("Node not found.")
}

// Retrieve an object previously serialized in JSON and remove it from the storage.
func (c *Client) GetAndRemove(key string, data interface{}) error {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	resp := &model.OvoResponse{Data: &model.OvoKVResponse{}}
	if s != nil {
		rs, err := s.Get(createGetAndRemoveEndpoint(s.node.Host, s.port, key), nil, resp, nil)
		if err != nil {
			// try get data from the twins
			done := true
			found := true
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					rs, errt := st.Get(createGetAndRemoveEndpoint(st.node.Host, st.port, key), nil, resp, nil)
					done = done && (errt ==nil)
					if errt == nil {
						if rs.status == 200 {
							errj := json.Unmarshal(resp.Data.(*model.OvoKVResponse).Data, data)
							found = found && (errj == nil)
						}
					}
					
				}
			}
			c.checkCluster()
			if done && found {
				return nil
			} else if done && !found {
				return errors.New("Key not found.")
			} else {
				return err
			}
		}
		if rs.status == 200 {
			err = json.Unmarshal(resp.Data.(*model.OvoKVResponse).Data, data)
			return err
		} else if rs.status == 404 {
			return errors.New("Key not found.")
		}
		return errors.New("Invalid data.")
	}
	return errors.New("Node not found.")
}

// Update an object with the newData if the oldData is equal to the stored data.
func (c *Client) UpdateValueIfEqual(key string, oldData interface{}, newData interface{}) error {
	hash := GetPositiveHashCode(key, maxServer)
	s := c.getSessionFromHash(hash)
	bOldData, err := json.Marshal(oldData)
	if err != nil {
		return err
	}
	bNewData, err := json.Marshal(newData)
	if err != nil {
		return err
	}
	mdata := &model.OvoKVUpdateRequest{Key: key, Data: bOldData, Hash: hash, NewData:bNewData}
	resp := &model.OvoResponse{}
	if s != nil {
		rs, err := s.Post(createUpdateValueIfEqualEndpoint(s.node.Host, s.port, key), mdata, resp, nil)
		if err != nil {
			// try get data from the twins
			done := true
			found := true
			for _, nd := range c.topology.GetTwins(s.node.Twins) {
				if st, ok := c.clients[nd.Name]; ok {
					rs, errt := st.Post(createUpdateValueIfEqualEndpoint(st.node.Host, st.port, key), mdata, resp, nil)
					done = done && (errt == nil)
					if errt == nil {
							found = found && (rs.status == 200)
					}
				}
			}
			c.checkCluster()
			if done && found {
				return nil
			} else if done && !found {
				return errors.New("Key not found or value not equal.")
			} else {
				return err
			}
		}
		if rs.status == 200 {
			return nil
		} else if rs.status == 403 {
			return errors.New("Forbidden operation: old value is not equal to the stored value.")
		} else if rs.status == 404 {
			return errors.New("Key not found.")
		}
		return errors.New("Invalid data.")
	}
	return errors.New("Node not found.")
}

