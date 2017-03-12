package main

import (
	"github.com/reactivex/rxgo/iterable"
	"fmt"
	
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/handlers"
	"reflect"
)

func main() {
	runGroupBy()
}

func runGroupBy() {
	type KV struct {
		K interface{}
		V int
	}

	items := []interface{}{KV{1,1}, KV{0,2}, KV{0,2}, KV{1,1}, KV{1,3}}
	it, err := iterable.New(items)
	if err != nil {
		fmt.Println(err)
	}

	stream1 := observable.From(it)

	key := func(item interface{}) interface{} {
		if i, ok := item.(KV); ok {
			return i.K
		}		
		return -1
	}

	stream2 := stream1.GroupBy(key)

	kvs := make(map[interface{}][]KV)
	onNext1 := handlers.NextFunc(func(item interface{}) {
		fmt.Println("Child onNext called: ",item)
		if p, ok := item.(KV); ok {
			kvs[p.K] = append(kvs[p.K],p)
		}
	})

	onNext := handlers.NextFunc(func(item interface{}) {
		if ob, ok := item.(observable.WrapObservable); ok {
			fmt.Println("Parent onNext called")
			s := ob.O.Subscribe(onNext1)
			go func() {
				<-s
			}()
		} else {
			fmt.Println("Parent onNext called, but not an observable. Got type: ", reflect.TypeOf(item))
		}
	})

	sub := stream2.Subscribe(onNext)
	<-sub
	fmt.Println(kvs)
}