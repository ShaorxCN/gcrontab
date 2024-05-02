package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DbConfig struct {
	Host     string `json:"db_host" env:"DB_HOST"`
	Port     string `json:"db_port" env:"DB_PORT"`
	User     string `json:"db_user" env:"DB_USER"`
	Name     string `json:"db_name" env:"DB_NAME"`
	Password string `json:"db_password" env:"DB_PASSWORD"`
}

func (*DbConfig) Init() error {
	return nil
}

func (*DbConfig) Restart() error {
	return nil
}

func (*DbConfig) Stop() {
}

func Test_getCfgFromENV(t *testing.T) {
	DefaultManagement.Register("db_config", &DbConfig{})
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "1234")
	os.Setenv("DB_USER", "evan")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_PASSWORD", "test123")
	ok := DefaultManagement.getConfigFromENV()
	if !ok {
		t.Errorf("test get environment config failed")
	}
	as := assert.New(t)
	dbConfig := DefaultManagement.Configs["db_config"].(*DbConfig)
	fmt.Println(dbConfig)
	as.Equal(dbConfig.Name, "test_db")
	as.Equal(dbConfig.Password, "test123")
}

func Test_withDefaultInfo(t *testing.T) {
	DefaultManagement.Register("db_config", &DbConfig{})
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "evan")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_PASSWORD", "test123")
	ok := DefaultManagement.getConfigFromENV()
	if !ok {
		t.Errorf("test get environment config failed")
	}
	as := assert.New(t)
	dbConfig := DefaultManagement.Configs["db_config"].(*DbConfig)
	fmt.Println(dbConfig)
	as.Equal(dbConfig.Port, "")
	as.Equal(dbConfig.Password, "test123")
}
