package husky

import (
	"strings"

	"github.com/spf13/viper"
)

type Config[T any] interface {
	Path(path string) Config[T]
	Name(name string) Config[T]
	Type(t string) Config[T]
	EnvPrefix(prefix string) Config[T]
	Load() *T
}

type _Config[T any] struct {
	v *viper.Viper
}

func (c *_Config[T]) Path(path string) Config[T] {
	c.v.AddConfigPath(path)
	return c
}

func (c *_Config[T]) Name(name string) Config[T] {
	c.v.SetConfigName(name)
	return c
}

func (c *_Config[T]) Type(t string) Config[T] {
	c.v.SetConfigType(t)
	return c
}

func (c *_Config[T]) EnvPrefix(prefix string) Config[T] {
	c.v.SetEnvPrefix(prefix)
	return c
}

func (c *_Config[T]) Load() *T {
	c.v.AutomaticEnv()
	c.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := c.v.ReadInConfig(); err != nil {
		panic(err)
	}
	var t T
	if err := c.v.Unmarshal(&t); err != nil {
		panic(err)
	}
	return &t
}

func NewConfig[T any]() Config[T] {
	v := viper.New()
	v.AddConfigPath("./")
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.SetEnvPrefix("APP")
	return &_Config[T]{v}
}
