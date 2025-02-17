// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

// Package globalconfig stores configuration which applies globally to both the tracer
// and integrations.
package globalconfig

import (
	"math"
	"sync"

	"github.com/google/uuid"
)

var cfg = &config{
	analyticsRate: math.NaN(),
	runtimeID:     uuid.New().String(),
	headersAsTags: make(map[string]string),
}

type config struct {
	mu            sync.RWMutex
	analyticsRate float64
	serviceName   string
	runtimeID     string
	headersAsTags map[string]string
}

// AnalyticsRate returns the sampling rate at which events should be marked. It uses
// synchronizing mechanisms, meaning that for optimal performance it's best to read it
// once and store it.
func AnalyticsRate() float64 {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()
	return cfg.analyticsRate
}

// SetAnalyticsRate sets the given event sampling rate globally.
func SetAnalyticsRate(rate float64) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.analyticsRate = rate
}

// ServiceName returns the default service name used by non-client integrations such as servers and frameworks.
func ServiceName() string {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()
	return cfg.serviceName
}

// SetServiceName sets the global service name set for this application.
func SetServiceName(name string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.serviceName = name
}

// RuntimeID returns this process's unique runtime id.
func RuntimeID() string {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()
	return cfg.runtimeID
}

// HeaderTag returns the tag assigned to the given header
func HeaderTag(header string) (tag string, ok bool) {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()
	tag, ok = cfg.headersAsTags[header]
	return tag, ok
}

// SetHeaderTag adds config for header `from` with tag value `to`
func SetHeaderTag(from, to string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.headersAsTags[from] = to
}

// HeaderTagsLen returns the length of globalconfig's headersAsTags map, 0 for empty map
func HeaderTagsLen() int {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()
	return len(cfg.headersAsTags)
}

// ClearHeaderTags assigns headersAsTags to a new, empty map
// It is invoked when WithHeaderTags is called, in order to overwrite the config
func ClearHeaderTags() {
	cfg.headersAsTags = make(map[string]string)
}
