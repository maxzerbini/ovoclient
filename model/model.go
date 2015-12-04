package model

const (
	MaxNodeNumber = 128
	Active = "ACTIVE"
	Inactive= "INACTIVE"
)

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
	Hash int32
}

type OvoKVUpdateRequest struct {
	Key string
	NewKey string
	Data []byte
	NewData []byte
	Hash int32
	NewHash int32
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
	Twins []string
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

// Get the active twin nodes
func (t *OvoTopology) GetTwins(names []string)(nodes []*OvoTopologyNode){
	nodes = make([]*OvoTopologyNode,0)
	for _,nd := range t.Nodes {
		for _,s := range names {
			if nd.Name == s {
				nodes = append(nodes, nd)
			}
		}
	}
	return nodes
}