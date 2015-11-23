package model

type OvoResponse struct {
	Status string
	Code string
	Data interface{}
}

type OvoKVRequest struct {
	Key string
	Data []byte
	Collection string
	TTL int
	Hash int
}

type OvoKVUpdateRequest struct {
	Key string
	NewKey string
	Data []byte
	NewData []byte
	Hash int
	NewHash int
}

type OvoKVResponse struct {
	Key string
	Data []byte
}

type OvoKVKeys struct {
	Keys []string
}

type OvoTopologyNode struct {
	Name string
	HashRange []int
	Host string
	Port int
	State string
}

type OvoTopology struct {
	Nodes []*OvoTopologyNode
}

type OvoResponseTopology struct {
	Status string
	Code string
	Data OvoTopology
}

type OvoResponseTopologyNode struct {
	Status string
	Code string
	Data OvoTopologyNode
}