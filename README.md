# ovo-client
The OVO Go Client Library.

## Latest Build
[![Build Status](https://drone.io/github.com/maxzerbini/ovoclient/status.png)](https://drone.io/github.com/maxzerbini/ovoclient/latest)

## Project Status
The project is completed, all the OVO API calls are implemented.

## Basic usage

### Put and Get an object
```Go
	var client = NewClientFromConfigPath("config.json")
	var testObj = &BigTestObject{ Name: "Max", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, 10000, 10000)}
	var err = client.Put("myObject", testObj, 0) // objects are serialized in JSON
	if err != nil {
		// menage error ...
	}
	var testMyObj = &BigTestObject{}
	err = client.Get("myObject", testMyObj)
	if err != nil {
		// menage error ...
	}
```

### Put and Get a raw object
```Go
	var client = NewClientFromConfigPath("config.json")
	var myImage = []byte(`1ab53b1237bf87329fh923...`)
	var err = client.PutRawData("myImage", myImage, 0)
	if err != nil {
		// menage error ...
	}
	image, err := client.GetRawData("myImage") // return a []byte 
	if err != nil {
		// menage error ...
	}
```

### Increment a counter
```Go
    var client = NewClientFromConfigPath("config.json")
	var count, err = client.Increment("myCounter", 1, 0)
	if err != nil {
		// menage error ...
	}
	initVal, err := client.GetCounter("myCounter")
	if err != nil {
		// menage error ...
	} else {
		printf("The initial value of %s is %d\r\n", "myCounter", count)
	}
	for i := 0; i < 10; i++ {
		count, err = client.Increment("myCounter", 1, 0)
		if err != nil {
			// menage error ...
		} else {
			printf("The value of %s is %d\r\n", "myCounter", count)
		}
	}
	count, err = client.SetCounter("myCounter", 0, 0)
	if err != nil {
		// menage error ...
	} else {
		printf("The value of %s is resetted to %d\r\n", "myCounter", count)
	}
```

## Client configuration

### Creating the client
There are three functions that can be used to configure the OVO client:
- using the default constructor function _NewClient()_ that creates the client reading the configuration file from the default directory "./config.json"
- using the constructor function _NewClientFromConfigPath(configpath string)_, passing a path and file name of the JSON configuration file
- using the constructor function _NewClientFromConfig(config *Configuration)_, passing a valid Configuration object in input

### The configuration file
The config.json file has this format
```JSON
{
	"ClusterNodes":
	[
		{"Host":"localhost","Port":"5050"}
	]
}
```
It can contain a list of one or more OVO node.

## Acknowledgments
I am indebted to Jason McVetta and his useful REST and HTTP client [Napping](https://github.com/jmcvetta/napping).
