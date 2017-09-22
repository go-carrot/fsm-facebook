package fsmfacebook

import "errors"

type CachedTraverser struct {
	uuid         string
	currentState string
	Data         map[string]interface{}
}

func (c *CachedTraverser) UUID() string {
	return c.uuid
}

func (c *CachedTraverser) SetUUID(newUUID string) {
	c.uuid = newUUID
}

func (c *CachedTraverser) CurrentState() string {
	return c.currentState
}

func (c *CachedTraverser) SetCurrentState(newState string) {
	c.currentState = newState
}

func (c *CachedTraverser) Upsert(key string, value interface{}) error {
	c.Data[key] = value
	return nil
}

func (c *CachedTraverser) Fetch(key string) (interface{}, error) {
	if val, ok := c.Data[key]; ok {
		return val, nil
	}
	return nil, errors.New("Key `" + key + "` is not set")
}

func (c *CachedTraverser) Delete(key string) error {
	delete(c.Data, key)
	return nil
}
