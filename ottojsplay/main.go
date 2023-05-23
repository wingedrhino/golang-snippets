// This program attempts to evaluate arbitrary JavaScript via Otto and at the
// same time determine a good number at which to parallelize its execution.
// The idea is to allow a mechanism where small snippets of buiness logic can be
// defined as JavaScript and thus evaluated both on the browser and on the
// backend server.
package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

// ToPretty returns an indented JSON representation of the interface. If there
// is an error, it uses fmt.Sprintf instead of json.MarshalIndent
func ToPretty(i interface{}) string {
	bytes, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		str := fmt.Sprintf("%+v", i)
		return str
	}
	strVal := string(bytes)
	return strVal
}

// MustToJSON returns a map[string]interface{} from a string that it assumes to
// be valid JSON. If the said string is not a valid JSON, it errors out.
func MustToJSON(str string) map[string]interface{} {
	res := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &res)
	if err != nil {
		log.Fatalf("Error unmarshaling string from JSON text to map[string]interface{}; error: %+v\n", err)
	}
	return res
}

// EvalOtto is a simple expression evaluator that uses Otto.
// The first argument is a `string` containing some JavaScript statements, which
// are enclosed into a function body and executed. There must thus be a return
// statement in there if you want to retrieve a non-null value as a result.
// The second argument is a `map[string]interface{}` that is injected into the
// program's environemnt as a variable named `args`. The function is capable of
// accessing the contents of args.
func EvalOtto(statements string, args map[string]interface{}) (interface{}, error) {
	vm := otto.New()
	err := vm.Set("args", args)
	if err != nil {
		return nil, err
	}
	program := fmt.Sprintf("function doSomething(){\n%s\n}\ndoSomething()", statements)
	vmValue, err := vm.Run(program)
	if err != nil {
		return vmValue, err
	}
	value, err := vmValue.Export()
	if err != nil {
		err = errors.Wrapf(err, "error calling value.Export()")
		return nil, err
	}
	return value, nil
}

// Executable holds arguments needed to run a program via OttoEval
type Executable struct {
	Program string                 `json:"program"`
	Args    map[string]interface{} `json:"args"`
}

// Eval evaluates this Executable via Otto
func (e Executable) Eval() (interface{}, error) {
	return EvalOtto(e.Program, e.Args)
}

// GetExecutable returns a random executable task from a list of executable
// tasks
func GetExecutable() Executable {
	exes := []Executable{
		{
			Program: `
			var sum = 0;
			if (args.some_bool) {
				for (var i = 0; i < args.numbers.length; i++) {
					sum = sum + args.numbers[i];
				}
			}
			return sum;
		`,
			Args: MustToJSON(`
				{
					"numbers": [1,2,3,4,5,6,7,8,9,10],
					"some_string": "sdsdsd",
					"some_bool": true
				}`),
		},
		{
			Program: `
			var sum = 0;
			if (args.some_bool) {
				for (var i = 0; i < args.numbers.length; i++) {
					sum = sum + args.numbers[i];
				}
			}
			// last value here
			return sum;
		`,
			Args: MustToJSON(`
				{
					"numbers": [1,2,3,4,5,6,7,8,9,10],
					"some_string": "sdsdsd",
					"some_bool": true
				}`),
		},
		{
			Program: `
			var sum = 0;
			if (args.some_bool) {
				for (var i = 0; i < args.numbers.length; i++) {
					sum = sum + args.numbers[i];
				}
			}
			// test comment
			return sum;
		`,
			Args: MustToJSON(`
				{
					"numbers": [1,2,3,4,5,6,7,8,9,10],
					"some_string": "sdsdsd",
					"some_bool": true
				}`),
		},
		{
			Program: `return 235;`,
			Args:    nil,
		},
		{
			Program: `return 235; // some value to return`,
			Args:    nil,
		},
		{
			Program: `return 236; // some other value to return`,
			Args:    nil,
		},
		{
			Program: `return args.x + args.y;`,
			Args:    MustToJSON(`{"x": 25, "y":2333}`),
		},
		{
			Program: `return args.y + args.x;`,
			Args:    MustToJSON(`{"x": 255, "y":2333}`),
		},
		{
			Program: `return args.x + args.y;`,
			Args:    MustToJSON(`{"x": 25, "y":23}`),
		},
	}
	i := GetRand(int64(len(exes)))
	return exes[i]
}

// GetRand returns a random int64 thats < max. It panics if there is an error
// fetching a random number
func GetRand(max int64) int64 {
	i, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Fatalf("Error generating random number: %+v\n", err)
	}
	return i.Int64()
}

// doOnce executes a single task
func doOnce(t Executable, wg *sync.WaitGroup) {
	_, err := t.Eval()
	if err != nil {
		log.Fatalf("Error processing task: %+v\n", err)
	}
	wg.Done()
}

// process processes a list of tasks read from a channel
func process(tasks *chan Executable, count int, done *chan bool) {
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		t := <-*tasks
		go doOnce(t, &wg)
	}
	wg.Wait()
	*done <- true
}

func main() {
	count := flag.Int("count", 40000, "number of iterations to run program for")
	maxGoroutines := flag.Int("routines", 0, "number of concurrent worker threads to run")
	flag.Parse()
	goMaxProcs := runtime.GOMAXPROCS(0)
	if *maxGoroutines == 0 {
		*maxGoroutines = 2 * goMaxProcs
	}
	done := make(chan bool)
	tasks := make(chan Executable, *maxGoroutines)
	tStart := time.Now()
	log.Printf("Operation Count: %d\n", *count)
	log.Printf("Max Goroutine Count: %d\n", *maxGoroutines)
	log.Printf("GOMAXPROCS: %d\n", goMaxProcs)
	log.Println("Begin Operation")
	go process(&tasks, *count, &done)
	for i := 0; i < *count; i++ {
		task := GetExecutable()
		tasks <- task
	}
	<-done
	tEnd := time.Now()
	tDiff := tEnd.Sub(tStart)
	log.Println("End Operation")
	log.Printf("Total time: %s\n", tDiff.String())
	log.Printf("Average time: %s\n", (tDiff / time.Duration(*count)).String())
}
