# FrontRunner
Frontrunner: Execute functions concurrently and return the result of the first to finish.

## How to use

### Create a new Runner
The Runner struct in the runner package encapsulates all required functionalities for task execution.

* Initialising without tasks
    ```
    r := runner.NewRunner[int]()
    ```
* Initialising with tasks
    ```
    r := runner.NewRunner(
        func() int { return 1 },
        func() int { return 1 * 5 },
    )
    ```

### Add tasks to Runner
Tasks can be added to runner after creation
```
r.Add(
    func() int { return 2 },
    func() int { return 2 * 5 },
)
r.Add(func() int { return 3 })
```

### Get first result of tasks
```
var res int
res, err := runner.First()
```

### Get first k results of tasks
```
var res []int
k := 2
res, err := runner.FirstK(k)
```

### Get first result with timeout condition
```
var (
    res int
    ok bool
)
res, ok, err := runner.FirstWithTimeout(time.Second)
```

### Get first k results with timeout condition
```
var (
    res []int
    ok bool
)
k := 2
res, ok, err := runner.FirstKWithTimeout(k, time.Second)
```

## Thread safety
Each operation on Runner struct is locked by Mutex, this means that each method call of a single runner will be synchronous, with only one critical operation running at a time.
```
// r.mu sync.Mutex

r.mu.Lock()
defer r.mu.Unlock()
```