package result

import "sync"

type (
	Result struct {
		sync.RWMutex
		query  string
		result [][]string
	}
)

func NewResult() *Result {
	return &Result{
		result: make([][]string, 0),
	}
}

func (r *Result) AddResult(rst [][]string) {
	r.Lock()
	defer r.Unlock()

	r.result = append(r.result, rst...)
}

func (r *Result) AddResult2(rst []string) {
	r.Lock()
	defer r.Unlock()

	r2 := make([][]string, 1)
	r2[0] = rst

	r.result = r2
}

func (r *Result) HasResult() bool {
	r.RLock()
	defer r.RUnlock()

	return len(r.result) > 0
}

func (r *Result) AddQuery(q string) {
	r.Lock()
	defer r.Unlock()

	r.query = q
}

func (r *Result) GetQuery() string {
	r.RLock()
	defer r.RUnlock()

	return r.query
}

func (r *Result) GetResult() chan []string {
	r.Lock()

	out := make(chan []string)

	go func() {
		defer close(out)
		defer r.Unlock()

		for _, r := range r.result {
			out <- r
		}
	}()

	return out
}
