package etcd

import (
	"config"
	"testing"
)

func Test_etcd(t *testing.T) {
	// Load config first
	config.LoadConfig()

	_, err := KeysAPI()
	if err != nil {
		t.Error(err)
	}
}

func Test_CreateDir(t *testing.T) {
	// Load config first
	config.LoadConfig()

	Kapi, err := KeysAPI()
	if err != nil {
		t.Error(err)
	}
	// CreateDir
	err = Kapi.CreateDir("/dddir")
	if err != nil {
		t.Error(err)
	}
}

func Test_DeleteDir(t *testing.T) {
	// Load config first
	config.LoadConfig()

	Kapi, err := KeysAPI()
	if err != nil {
		t.Error(err)
	}
	// DeleteDir
	err = Kapi.DeleteDir("/dddir")
	if err != nil {
		t.Error(err)
	}
}

func Test_SetValue(t *testing.T) {
	// Load config first
	config.LoadConfig()

	Kapi, err := KeysAPI()
	if err != nil {
		t.Error(err)
	}
	// SetValue
	err = Kapi.SetValue("/test_key", "value", 0)
	if err != nil {
		t.Error(err)
	}
}

func Test_GetValue(t *testing.T) {
	// Load config first
	config.LoadConfig()

	Kapi, err := KeysAPI()
	if err != nil {
		t.Error(err)
	}
	// GetValue
	value, err := Kapi.GetValue("/test_key")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(value)
	}
}

func Test_DeleteKey(t *testing.T) {
	// Load config first
	config.LoadConfig()

	Kapi, err := KeysAPI()
	if err != nil {
		t.Error(err)
	}
	// DeleteKey
	err = Kapi.DeleteKey("/test_key")
	if err != nil {
		t.Error(err)
	}
}

func Test_CreateInOrder(t *testing.T) {
	// Load config first
	config.LoadConfig()

	Kapi, err := KeysAPI()
	if err != nil {
		t.Error(err)
	}

	dir := "/a/b/c/d"
	// CreateInOrder
	err = Kapi.CreateInOrder(dir)
	if err != nil {
		t.Error(err)
	}

	Kapi.DeleteDir(dir)
}
