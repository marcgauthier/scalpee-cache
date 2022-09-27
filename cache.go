package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/* each item in the cache contain a float and the time when it will expire
 */
type cacheItem struct {
	value      float64
	expireTime int64
}

/* each cache contains a map that contains all the cache items and
   the default cache expired time. and it's own Mutex.
*/
type CacheObj struct {
	items                map[string]cacheItem
	defaultkeepItemsMins int
	mutex                sync.RWMutex
}

/* create a new cache initiate the memory and set the default value
   for the expiry time
*/
func New(keepItemsMins int) *CacheObj {

	c := new(CacheObj)
	c.items = make(map[string]cacheItem)
	c.defaultkeepItemsMins = keepItemsMins
	go c.start()

	return c
}

/* set a value in the cache, if a value exist it is overwritten. if the TTL
   is set to 0 then the default value is used.  Otherwise TTL contain the
   amount of minutes to keep data.
*/

func (c *CacheObj) Set(key string, v float64, ttl int) {
	var expire int64
	c.mutex.Lock()
	if ttl == 0 {
		expire = time.Now().Add(time.Second * 60 * time.Duration(c.defaultkeepItemsMins)).Unix()
	} else {
		expire = time.Now().Add(time.Second * 60 * time.Duration(ttl)).Unix()
	}
	// delete existing cache
	delete(c.items, key)
	c.items[key] = cacheItem{value: v, expireTime: expire}
	c.mutex.Unlock()
}

/* get a key from the cache
 */

func (c *CacheObj) Get(key string) (float64, error) {

	var result cacheItem
	var value float64
	var ok bool

	c.mutex.RLock()
	if result, ok = c.items[key]; !ok {
		c.mutex.RUnlock()
		return 0, errors.New("not found")
	}

	if result.expireTime < time.Now().Unix() {
		c.mutex.RUnlock()
		return 0, errors.New("expired")
	}

	value = result.value
	c.mutex.RUnlock()
	return value, nil
}

/* delete expired cache
 */
func (c *CacheObj) start() {

	fmt.Println("starting cache")
	for {
		time.Sleep(time.Second)
		c.mutex.Lock()
		for k, i := range c.items {
			if i.expireTime < time.Now().Unix() {
				//fmt.Println("removing " + k)
				delete(c.items, k)
			}
		}
		c.mutex.Unlock()
	}

}
