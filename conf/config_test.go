package conf_test

import (
	"os"
	"testing"

	"github.com/Murphychih/cmdb/conf"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromToml(t *testing.T) {
	should := assert.New(t)
	err := conf.LoadConfigFromToml("../etc/demo.toml")
	if should.NoError(err) {
		global := conf.LoadGloabal()
		should.Equal("demo", global.App.Name)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	should := assert.New(t)
	os.Setenv("MYSQL_DATABASE", "unit_test")

	err := conf.LoadConfigFromEnv()
	if should.NoError(err) {
		global := conf.LoadGloabal()
		should.Equal("unit_test", global.MySQL.Database)
	}
}

func TestGetDB(t *testing.T) {
	should := assert.New(t)
	err := conf.LoadConfigFromToml("../etc/demo.toml")
	if should.NoError(err) {
		global := conf.LoadGloabal()
		global.MySQL.GetDB()
	}
}
