package ovoclient

import (
	"strconv"
	"testing"
	"time"
)

var client = NewClientFromConfigPath("config.json")

const (
	MaxBigObjectItems = 100 // 1000
	BigObjectSize     = 1000
)

func init() {
	LogEnabled = true
}

type TestObject struct {
	Name      string
	Surname   string
	BirthDate time.Time
	Id        int32
}

type BigTestObject struct {
	Name      string
	Surname   string
	BirthDate time.Time
	Id        int32
	LotOfData []byte
}

/**/

func TestConfigurationLoad(t *testing.T) {
	client := NewClientFromConfigPath("config.json")
	if client == nil {
		t.Fail()
	}
}

func TestPutRawData(t *testing.T) {
	var jsonStr = []byte(`{"Key":"data3","Data":"dGVzdA==","TTL":103}`)
	var err = client.PutRawData("test123", jsonStr, 0)
	if err != nil {
		t.Fail()
	} else {
		t.Logf("Test PutRawData done")
	}
}

func TestGetRawData(t *testing.T) {
	var jsonStr = []byte(`{"Key":"data3","Data":"dGVzdA==","TTL":103}`)
	var err = client.PutRawData("test1234", jsonStr, 0)
	if err != nil {
		t.Fail()
	} else {
		resp2, err := client.GetRawData("test1234")
		if err != nil {
			t.Logf("Error: %v", err)
			t.Fail()
		}
		t.Logf("result = %v", resp2)
	}

}

func TestPutTestObject(t *testing.T) {
	var testObj = &TestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111}
	var err = client.Put("testobj555", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		t.Logf("Test PutTestObject done")
	}
}

func TestGetTestObject(t *testing.T) {
	var testObj = &TestObject{}
	var err = client.Get("testobj555", testObj)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Fail()
	} else {
		t.Logf("result = %v", *testObj)
	}
}

func TestNotFound(t *testing.T) {
	var testObj = &TestObject{}
	var key = "notfound"
	var err = client.Get(key, testObj)
	if err != nil {
		t.Logf("Key not found %s :  %v", key, err)

	} else {
		t.Errorf("result = %v", *testObj)
		t.Fail()
	}
}

func TestPutBigObject(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, BigObjectSize, BigObjectSize)}
	var err = client.Put("bigobj1", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		t.Logf("Test PutBigObject done")
	}
}

func TestGetBigObject(t *testing.T) {
	var testObj = &BigTestObject{}
	var err = client.Get("bigobj1", testObj)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Fail()
	} else {
		t.Logf("result = %v", *testObj)
	}
}

func TestPutVeryBigObject(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, BigObjectSize*10, BigObjectSize*10)}
	var err = client.Put("bigobj2", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		t.Logf("Test PutVeryBigObject done")
	}
}

func TestGetVeryBigObject(t *testing.T) {
	var testObj = &BigTestObject{}
	var err = client.Get("bigobj2", testObj)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Fail()
	} else {
		t.Logf("result = %v", *testObj)
	}
}

func TestPutALotOfBigObject(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 0, LotOfData: make([]byte, 10000, 10000)}
	for i := 0; i < MaxBigObjectItems; i++ {
		testObj.Id = 100 + int32(i)
		var err = client.Put("bigobj3_"+strconv.Itoa(i), testObj, 0)
		if err != nil {
			t.Fail()
		}
	}
	t.Logf("Test PutALotOfBigObject done for %d objects", MaxBigObjectItems)
}
func TestGetALotOfBigObject(t *testing.T) {
	for i := 0; i < MaxBigObjectItems; i++ {
		var testObj = &BigTestObject{}
		var err = client.Get("bigobj3_"+strconv.Itoa(i), testObj)
		if err != nil {
			t.Logf("Error: %v", err)
			t.Fail()
		} else {
			t.Logf("testObj.Id = %v", testObj.Id)
		}
	}
}

/**/
/**/
func TestPutALotOfBigObject2(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, 100, 100)}
	for i := 0; i < MaxBigObjectItems; i++ {
		testObj.Id = 10000 + +int32(i)
		var err = client.Put("bigobj4_"+strconv.Itoa(i), testObj, 0)
		if err != nil {
			t.Fail()
		}
	}
	t.Logf("Test PutALotOfBigObject2 done for %d objects", MaxBigObjectItems)
}

/**/

func TestCount(t *testing.T) {
	var count = client.Count()
	if len(count) == 0 {
		t.Fail()
	} else {
		t.Logf("Counters = %v", count)
	}
}

func TestKeys(t *testing.T) {
	var keys = client.Keys()
	if len(keys) == 0 {
		t.Fail()
	} else {
		t.Logf("Keys number = %d\r\n", len(keys))
	}
}

func TestDelete(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, BigObjectSize, BigObjectSize)}
	var err = client.Put("bigobjToDelete", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		err := client.Delete("bigobjToDelete")
		if err != nil {
			t.Fail()
		} else {
			result := &BigTestObject{}
			err := client.Get("bigobjToDelete", result)
			if result != nil && err == nil {
				t.Fail()
			} else {
				t.Log("Object deleted.\r\n")
			}
		}

	}
}

func TestGetAndRemove(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, BigObjectSize, BigObjectSize)}
	var err = client.Put("bigobjToRemove", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		var testObjRemoved = &BigTestObject{}
		err := client.GetAndRemove("bigobjToRemove", testObjRemoved)
		if err != nil && testObjRemoved == nil {
			t.Fail()
		} else {
			t.Log("Removed object: %v", testObjRemoved)
			result := &BigTestObject{}
			err := client.Get("bigobjToRemove", result)
			if result != nil && err == nil {
				t.Fail()
			} else {
				t.Log("Object removed.\r\n")
			}
		}

	}
}

func TestUpdateValueIfEqual(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, BigObjectSize, BigObjectSize)}
	var testNewObj = &BigTestObject{Name: "Max", Surname: "Zerbini", BirthDate: time.Now(), Id: 112, LotOfData: make([]byte, BigObjectSize, BigObjectSize)}
	var err = client.Put("bigobjToUpdate", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		err := client.UpdateValueIfEqual("bigobjToUpdate", testObj, testNewObj)
		if err != nil {
			t.Fail()
		} else {
			result := &BigTestObject{}
			err := client.Get("bigobjToUpdate", result)
			if err != nil {
				t.Fail()
			} else {
				if result.Id == testObj.Id {
					t.Fail()
				} else {
					t.Logf("Updated ojbect is %v\r\n", *result)
				}
			}
		}

	}
}

func TestIncrement(t *testing.T) {

	var count, err = client.Increment("myCounter", 1, 0)
	if err != nil {
		t.Fail()
	}
	initVal, err := client.GetCounter("myCounter")
	if err != nil {
		t.Fail()
	} else {
		t.Logf("The initial value of %s is %d\r\n", "myCounter", count)
	}
	for i := 0; i < 10; i++ {
		count, err = client.Increment("myCounter", 1, 0)
	}
	if count < (initVal + 10) {
		t.Fail()
	} else {
		t.Logf("The value of %s is %d\r\n", "myCounter", count)
	}
	count, err = client.SetCounter("myCounter", 0, 0)
	if count > 0 {
		t.Fail()
	} else {
		t.Logf("The value of %s is resetted to %d\r\n", "myCounter", count)
	}
}
