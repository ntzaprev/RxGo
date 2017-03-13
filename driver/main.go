package main

import (
	"github.com/reactivex/rxgo/iterable"
	"fmt"
	
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/handlers"
	"reflect"
	//"time"
	"github.com/reactivex/rxgo/fx"
)

func main() {
	//runGroupBy()
	//runCombineLatest()
	runZip()
}

func runZip() {
	items1 := []interface{}{0, 1, 2, 3, 4, 5}

	it1, err := iterable.New(items1)
	if err != nil {
		fmt.Println(err)
	}

	stream1 := observable.From(it1)
	items2 := []interface{}{0, 1, 2, 3, 4, 5,6,7,8}
	
	it2, err := iterable.New(items2)
	if err != nil {
		fmt.Println(err)
	}

	stream2 := observable.From(it2)

	/*fin := make(chan struct{})
	stream2 := observable.Interval(fin, 1000*time.Millisecond)
	OnNext2 := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			if num >= 5 {
				fin <- struct{}{}
				close(fin)
			}
		}
	})

	sub:= stream2.Subscribe(OnNext2)
	go func(){
		<-sub
	}()*/

	comb := fx.CombinableFunc(func(in1,in2 interface{})interface{} {
		var r1,r2 int
		if i,ok := in1.(int); ok {
			r1 = i
		}
		if i,ok := in2.(int); ok {
			r2 = i
		}
		return fmt.Sprintf("[%v,%v]",r1,r2)
	})
	stream3 := observable.Zip(stream1, stream2, comb)

	var pairs []string
	OnNext := handlers.NextFunc(func(item interface{}){
		fmt.Println("OnNext on combined callded")
		if s,ok := item.(string); ok {
			pairs = append(pairs,s)
		}
	})
	s := stream3.Subscribe(OnNext)
	<-s
	fmt.Println(pairs)
}

func runCombineLatest() {
	items1 := []interface{}{0, 1, 2, 3, 4, 5}

	it1, err := iterable.New(items1)
	if err != nil {
		fmt.Println(err)
	}

	stream1 := observable.From(it1)
	items2 := []interface{}{0, 1, 2, 3, 4, 5}
	
	it2, err := iterable.New(items2)
	if err != nil {
		fmt.Println(err)
	}

	stream2 := observable.From(it2)

	/*fin := make(chan struct{})
	stream2 := observable.Interval(fin, 1000*time.Millisecond)
	OnNext2 := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			if num >= 5 {
				fin <- struct{}{}
				close(fin)
			}
		}
	})

	sub:= stream2.Subscribe(OnNext2)
	go func(){
		<-sub
	}()*/

	comb := fx.CombinableFunc(func(in1,in2 interface{})interface{} {
		var r1,r2 int
		if i,ok := in1.(int); ok {
			r1 = i
		}
		if i,ok := in2.(int); ok {
			r2 = i
		}
		return fmt.Sprintf("[%v,%v]",r1,r2)
	})
	stream3 := observable.CombineLatest(stream1, stream2, comb)

	var pairs []string
	OnNext := handlers.NextFunc(func(item interface{}){
		fmt.Println("OnNext on combined callded")
		if s,ok := item.(string); ok {
			pairs = append(pairs,s)
		}
	})
	s := stream3.Subscribe(OnNext)
	<-s
	fmt.Println(pairs)
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