package cache

import (
	"bytes"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	uuid "github.com/satori/go.uuid"
	"gocv.io/x/gocv"
)

// ImageCache .
type ImageCache struct {
	Client    *memcache.Client
	Namespace *uuid.UUID
}

// NewCache .
func NewCache(active bool, address string) (*ImageCache, error) {
	if !active {
		fmt.Println("============ Cache disabled ==============")
		return &ImageCache{
			Client:    nil,
			Namespace: nil,
		}, nil
	}

	mc := memcache.New(address)

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &ImageCache{
		Client:    mc,
		Namespace: &id,
	}, nil
}

// GetBytes .
func (c ImageCache) GetBytes(key string) ([]byte, error) {

	if c.Client == nil {
		return nil, nil
	}

	item, err := c.Client.Get(key)

	if err == memcache.ErrCacheMiss {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

// GetMat .
func (c ImageCache) GetMat(key string) (*gocv.Mat, error) {

	if c.Client == nil {
		return nil, nil
	}

	item, err := c.Client.Get(key)

	if err != nil {
		return nil, err
	}

	mat, err := gocv.IMDecode(item.Value, gocv.IMReadUnchanged)

	return &mat, err
}

// AddMat .
func (c ImageCache) AddMat(key string, mat *gocv.Mat) error {
	buf := mat.ToBytes()
	return c.AddBytes(key, buf)
}

// AddBytes .
func (c ImageCache) AddBytes(key string, buf []byte) error {

	if c.Client == nil {
		return nil
	}

	item := memcache.Item{
		Key:   key,
		Value: buf,
	}

	return c.Client.Add(&item)
}

// GenerateHash .
func (c ImageCache) GenerateHash(keys ...string) string {
	buf := bytes.Buffer{}
	for _, v := range keys {
		buf.WriteString(v)
	}

	return uuid.NewV5(*c.Namespace, buf.String()).String()
}
