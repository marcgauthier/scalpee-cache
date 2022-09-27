package cache

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

// test function
func TestCache(t *testing.T) {

	c := New(1)

	for i := 0; i < 5; i++ {
		c.Set(strconv.Itoa(i), float64(i), 0)
	}

	for i := 0; i < 5; i++ {

		v, err := c.Get(strconv.Itoa(i))

		if err != nil {
			t.Errorf("missing")
		}
		if v != float64(i) {
			t.Errorf("value not the same")
		}

	}

	fmt.Println("pausing for 30 secs")
	time.Sleep(time.Second * 30)
	c.Set("1", 12, 0)

	fmt.Println("pausing for 40 secs")
	time.Sleep(time.Second * 40)

	x, err := c.Get("2")
	if err != nil {
		// good item can't be found
	} else {
		t.Errorf("cache item didn't expired value is " + strconv.FormatFloat(x, 'f', 10, 64))
	}

	v, err := c.Get("1")
	if err != nil {
		t.Errorf("cache item shoult still exists")
	} else {
		if v != 12 {
			t.Errorf("item 1 was not changed to value=12")
		}

	}

}
