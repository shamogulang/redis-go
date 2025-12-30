package main

import "sync"

var SETs = map[string]string{}

var SETsMu = sync.RWMutex{}
