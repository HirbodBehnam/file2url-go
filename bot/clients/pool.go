package clients

import (
	"container/list"
	"context"
	"github.com/gotd/td/tg"
	"sync"
	"time"
)

// apiPool contains a pool for telegram api, so we don't create and delete them for each request
type apiPool struct {
	// A list of *ApiPoolElement
	l  *list.List
	mu sync.Mutex
}

// ApiPoolElement contains the client, the cancel function to close the api, and last access time
type ApiPoolElement struct {
	accessTime time.Time
	tgApi      *tg.Client
	cancel     context.CancelFunc
}

// API will return the telegram API
func (e *ApiPoolElement) API() *tg.Client {
	return e.tgApi
}

// newApiPool creates a new api pool
func newApiPool() *apiPool {
	pool := &apiPool{
		l: list.New(),
	}
	go pool.gc()
	return pool
}

func (p *apiPool) gc() {
	// gcCycle says the interval which we must check the old clients
	const gcCycle = time.Minute
	// What is the max age of client
	const clientTTL = time.Minute
	for {
		time.Sleep(gcCycle)
		p.mu.Lock()
		for e := p.l.Back(); e != nil; {
			element := e.Value.(*ApiPoolElement)
			if time.Since(element.accessTime) > clientTTL { // This api must be removed
				element.cancel() // at first cancel the api
				temp := e        // create a temp to...
				e = e.Prev()     // ...iterate the list
				p.l.Remove(temp) // remove from list
			} else {
				break // no point in looking other variables. Because of the insertion cycle they have lived less than this
			}
		}
		p.mu.Unlock()
	}
}

// Get gets a new client from pool
// If nothing exists, it will create a new client
func (p *apiPool) Get() (*ApiPoolElement, error) {
	// Check if anything exists
	p.mu.Lock()
	if e := p.l.Front(); e != nil {
		p.l.Remove(e)
		p.mu.Unlock()
		return e.Value.(*ApiPoolElement), nil
	}
	p.mu.Unlock()
	// We shall create a new client
	ctx, cancel := context.WithCancel(context.Background())
	client, err := instantiateClient(ctx)
	if err != nil {
		cancel()
		return nil, err
	}
	// Return the element
	return &ApiPoolElement{
		tgApi:  client,
		cancel: cancel,
	}, nil
}

func (p *apiPool) Put(e *ApiPoolElement) {
	e.accessTime = time.Now() // update the time
	p.mu.Lock()
	p.l.PushFront(e) // push to front. Next Get before Put will return this value
	p.mu.Unlock()
}
