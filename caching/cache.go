package caching

import (
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

var Cache = cache.New(10*time.Minute, 10*time.Minute)

// RateLimiter keeps track of IP-based rate limits using go-cache
type RateLimiter struct {
	cache *cache.Cache
	rate  rate.Limit
	burst int
	ttl   time.Duration // Time to live for each IP rate limit
}

// NewRateLimiter initializes a new RateLimiter with go-cache (in-memory cache)
func NewRateLimiter(r rate.Limit, b int, ttl time.Duration) *RateLimiter {
	c := cache.New(ttl, ttl)
	return &RateLimiter{
		cache: c,
		rate:  r,
		burst: b,
		ttl:   ttl,
	}
}

// GetCache retrieves the current request count for a specific IP from go-cache
func (rl *RateLimiter) GetCache(ip string) (int64, error) {
	limiterKey := fmt.Sprintf("limiter:%s", ip)
	val, found := rl.cache.Get(limiterKey)
	if !found {
		return 0, nil // No previous record means it's the first request
	}

	count, ok := val.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid value for limiter key")
	}

	return count, nil
}

// SetCache increments the request count for a specific IP in go-cache with TTL
func (rl *RateLimiter) SetCache(ip string) (int64, error) {
	limiterKey := fmt.Sprintf("limiter:%s", ip)

	// Increment the request count
	count, _ := rl.GetCache(ip)
	count++

	// Store the new count in the cache
	rl.cache.Set(limiterKey, count, rl.ttl)

	return count, nil
}

// getLimiter checks if the request should be throttled
func (rl *RateLimiter) getLimiter(ip string) (*rate.Limiter, error) {
	count, err := rl.GetCache(ip)
	if err != nil {
		return nil, err
	}
	fmt.Println(count)
	// If requests exceed the burst limit, return an error
	if count >= int64(rl.burst) {
		return nil, fmt.Errorf("rate limit exceeded for IP %s", ip)
	}

	// Increment the request count
	_, err = rl.SetCache(ip)
	if err != nil {
		return nil, err
	}

	// Return a rate limiter object
	return rate.NewLimiter(rl.rate, rl.burst), nil
}

// Middleware applies rate limiting to incoming requests
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		// Fetch or create a rate limiter for the IP address
		limiter, err := rl.getLimiter(ip)
		if err != nil {
			http.Error(w, "Too many requests, please slow down", http.StatusTooManyRequests)
			return
		}

		// Check if the rate limiter allows the request
		if !limiter.Allow() {
			http.Error(w, "Too many requests, please slow down", http.StatusTooManyRequests)
			return
		}

		// Proceed with the next handler if the request is allowed
		next.ServeHTTP(w, r)
	})
}

// cacheHandler handles caching for expensive operations
func cacheHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize go-cache for cache
	cache := cache.New(10*time.Minute, 10*time.Minute)

	key := "expensive_data"

	// Check if data is in cache
	value, found := cache.Get(key)
	if found {
		// Cache hit
		fmt.Fprintln(w, "Cache Hit: ", value)
	} else {
		// Cache miss, simulate expensive operation
		value = "Expensive Data Result"
		cache.Set(key, value, time.Minute)
		fmt.Fprintln(w, "Cache Miss: Data was fetched and cached:", value)
	}
}
