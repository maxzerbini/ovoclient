package ovoclient

import (
	"bytes"
)

func createTopologyEndpoint(host string, port string) string {
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(port)
	buffer.WriteString("/ovo/cluster")
	return buffer.String()
}

func createTopologyNodeEndpoint(host string, port string) string {
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(port)
	buffer.WriteString("/ovo/cluster/me")
	return buffer.String()
}

func createKeysEndpoint(host string, port string) string {
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(port)
	buffer.WriteString("/ovo/keys")
	return buffer.String()
}

func createKeyStorageEndpoint(host string, port string) string {
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(port)
	buffer.WriteString("/ovo/keystorage")
	return buffer.String()
}

func createGetKeyStorageEndpoint(host string, port string, key string) string {
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(port)
	buffer.WriteString("/ovo/keystorage/")
	buffer.WriteString(key)
	return buffer.String()
}

func createGetAndRemoveEndpoint(host string, port string, key string) string {
	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(host)
	buffer.WriteString(":")
	buffer.WriteString(port)
	buffer.WriteString("/ovo/keystorage/")
	buffer.WriteString(key)
	buffer.WriteString("/getandremove")
	return buffer.String()
}
