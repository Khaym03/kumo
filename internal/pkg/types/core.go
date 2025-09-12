package types

type Request struct {
	URL       string
	Metadata  map[string]any
	Collector string
}

type RequestQueue chan *Request

func (rq RequestQueue) Enqueue(reqs ...*Request) {
	for _, r := range reqs {
		rq <- r
	}
}
