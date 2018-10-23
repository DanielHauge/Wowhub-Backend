package Redis

import "reflect"

func ServeCacheAndUpdateBehind(key string, id int, expectedType interface{},fetcher func(id int, obj *interface{}) error) (chan interface{}, chan error) {
	channel := make(chan interface{})
	errorcheck := make(chan error)

	go func() {
		result := reflect.New(reflect.TypeOf(expectedType)).Interface()
		var e error
		if DoesKeyExist(key) {
			e = CacheGetResult(key, &result)
			if !CacheTimeout(key) { // If key is not in timeout, update cache.
				go func() {
					var Caching interface{}
					e = fetcher(id, &Caching)
					CacheSetResult(key, Caching)
				}()
			}
		} else {
			e = fetcher(id, &result)
			if e == nil {
				go CacheSetResult(key, result)
			}
		}
		if e != nil {
			errorcheck <- e
		} else {
			channel <- result
		}

	}()

	return channel, errorcheck

}
