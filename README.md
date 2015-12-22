# ovo-client
OVO Go Client Library

The project is completed, all the OVO API calls are implemented.

## Basic usage
```Go
	var client = NewClientFromConfigPath("config.json")
	var testObj = &BigTestObject{Name: "Max", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, 10000, 10000)}
	var err = client.Put("myObject", testObj, 0)
	if err != nil {
		panic
	}
	var testMyObj = &BigTestObject{}
	var err = client.Get("myObject", testMyObj)
```