package core

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

// Cache uint
var Cache = cache.New(24*time.Hour, 48*time.Hour)
