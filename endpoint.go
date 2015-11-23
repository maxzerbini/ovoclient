package ovoclient
import(
	"bytes"
)

func createTopologyEndpoint(host string, port string) string{
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(port)
	buffer.WriteString("/ovo/cluster")
	return buffer.String()
}

func createTopologyNodeEndpoint(host string, port string) string{
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(port)
	buffer.WriteString("/ovo/cluster/me")
	return buffer.String()
}