package ovoclient
import(
	"io/ioutil"
	"os"
	"encoding/json"
	"log"
)
type Node struct {
	Host string
	Port string
}

type Configuration struct {
	ClusterNodes []Node
}

func LoadConfiguration(path string) *Configuration {
	file, e := ioutil.ReadFile(path)
    if e != nil {
		log.Fatalf("Configuration file not found at %s", path)
        os.Exit(1)
    }
    var jsontype *Configuration = &Configuration{}
    json.Unmarshal(file, jsontype)
	return jsontype;
}