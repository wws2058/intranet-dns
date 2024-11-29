package cronjob

import (
	"fmt"
	"math/rand"
	"time"
)

var internalFunctionMaps = map[string]func() error{
	"test_function": testFunction,
}

func testFunction() error {
	rand.NewSource(time.Now().UnixNano())
	randomInt := rand.Intn(100)
	if randomInt%2 == 0 {
		return nil
	}
	return fmt.Errorf("%v is not even number", randomInt)
}
