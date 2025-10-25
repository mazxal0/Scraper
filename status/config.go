package status

import "sync"

var Statuses = map[string]*URLStatus{}

var mu sync.Mutex
