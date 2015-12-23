# ovo-client
OVO Go Client Library

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

### Put and Get an raw object
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