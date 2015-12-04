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
	var resp, err = client.PutRawData("test123", jsonStr, 0)
	if err != nil {
		t.Fail()
	} else {
		t.Logf("result = %v", *resp)
	}
}

func TestGetRawData(t *testing.T) {
	var jsonStr = []byte(`{"Key":"data3","Data":"dGVzdA==","TTL":103}`)
	var _, err = client.PutRawData("test123", jsonStr, 0)
	if err != nil {
		t.Fail()
	} else {
		resp2, err := client.GetRawData("test123")
		if err != nil {
			t.Logf("Error: %v", err)
			t.Fail()
		}
		t.Logf("result = %v", resp2)
	}

}

func TestPutTestObject(t *testing.T) {
	var testObj = &TestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111}
	var resp, err = client.Put("testobj555", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		t.Logf("result = %v", *resp)
	}
}

func TestGetTestObject(t *testing.T) {
	var testObj = &TestObject{}
	var resp, err = client.Get("testobj555", testObj)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Fail()
	} else {
		t.Logf("result = %v", *resp)
		t.Logf("result = %v", *testObj)
	}
}

func TestNotFound(t *testing.T) {
	var testObj = &TestObject{}
	var key = "notfound"
	var resp, err = client.Get(key, testObj)
	if err != nil {
		t.Logf("Key not found %s :  %v", key, err)

	} else {
		t.Errorf("result = %v", *resp)
		t.Errorf("result = %v", *testObj)
		t.Fail()
	}
}

func TestPutBigObject(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, BigObjectSize, BigObjectSize)}
	var resp, err = client.Put("bigobj1", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		t.Logf("result = %v", *resp)
	}
}

func TestGetBigObject(t *testing.T) {
	var testObj = &BigTestObject{}
	var resp, err = client.Get("bigobj1", testObj)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Fail()
	} else {
		t.Logf("result = %v", *resp)
		t.Logf("result = %v", *testObj)
	}
}

func TestPutVeryBigObject(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, BigObjectSize*10, BigObjectSize*10)}
	var resp, err = client.Put("bigobj2", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		t.Logf("result = %v", *resp)
	}
}

func TestGetVeryBigObject(t *testing.T) {
	var testObj = &BigTestObject{}
	var resp, err = client.Get("bigobj2", testObj)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Fail()
	} else {
		t.Logf("result = %v", *resp)
		t.Logf("result = %v", *testObj)
	}
}

func TestPutALotOfBigObject(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, 10000, 10000)}
	for i := 0; i < MaxBigObjectItems; i++ {
		var resp, err = client.Put("bigobj3_"+strconv.Itoa(i), testObj, 0)
		if err != nil {
			t.Fail()
		} else {
			t.Logf("result = %v", *resp)
		}
	}
}
func TestGetALotOfBigObject(t *testing.T) {
	for i := 0; i < MaxBigObjectItems; i++ {
		var testObj = &BigTestObject{}
		var resp, err = client.Get("bigobj3_"+strconv.Itoa(i), testObj)
		if err != nil {
			t.Logf("Error: %v", err)
			t.Fail()
		} else {
			t.Logf("result = %v", *resp)
		}
	}
}

/**/
/**/
func TestPutALotOfBigObject2(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, 100, 100)}
	for i := 0; i < MaxBigObjectItems; i++ {
		var resp, err = client.Put("bigobj4_"+strconv.Itoa(i), testObj, 0)
		if err != nil {
			t.Fail()
		} else {
			t.Logf("result = %v", *resp)
		}
	}
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
		t.Logf("Keys number = %d", len(keys))
	}
}

func TestDelete(t *testing.T) {
	var testObj = &BigTestObject{Name: "Massimo", Surname: "Zerbini", BirthDate: time.Now(), Id: 111, LotOfData: make([]byte, BigObjectSize, BigObjectSize)}
	var _, err = client.Put("bigobjToDelete", testObj, 0)
	if err != nil {
		t.Fail()
	} else {
		_, err := client.Delete("bigobjToDelete")
		if err != nil {
			t.Fail()
		} else {
			result := &BigTestObject{}
			_, err := client.Get("bigobjToDelete", result)
			if result != nil && err == nil {
				t.Fail()
			} else {
				t.Log("Object deleted.")
			}
		}

	}
}
