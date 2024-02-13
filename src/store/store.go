package main

import (
	"errors"
	"sync"
)

/*
This is meant to holds logic for the in-memory KV store
Features include:
	- Set Value
	- Get Value

Storage is of the mode 
{
	key1: [log1, log2, log3],
	key2: [log1a, log1b, log1c],
}
*/

type Store struct{
	mu sync.Mutex
	logs map[string][]string
}

// Create a new Store
func NewStore() *Store{
	return &Store{
		logs: make(map[string][]string),	
	}
}

// check if key exists
func (s *Store) keyExists(key string) (bool, []string){
	val, exists := s.logs[key]
	if !exists{
		return false, []string{}
	}
	return true, val
}

/* 
Add new value to a key
it creats the key if not exists
*/
func (s *Store)Append(key string, value string) *Store{
	s.mu.Lock()
	exists, logs := s.keyExists(key)
	// New Key not found
	if ! exists{
		s.logs[key] = append(s.logs[key], value)
		s.mu.Unlock()
		return s
	}
	// key exists
	s.logs[key] = append(logs, value)
	s.mu.Unlock()
	return s
}

// Get value of a key
// Returns error if key not found
func (s *Store)Get(key string) (string, error){
	s.mu.Lock()
	defer s.mu.Unlock()
	exists, logs := s.keyExists(key)
	if !exists{
		return "", errors.New("key not found")
	}
	// return the last item appended
	return logs[len(logs)-1], nil
}

// Get key value at an index
func (s *Store) GetAtIndex(key string, index int) (string, error){
	s.mu.Lock()
	defer s.mu.Unlock()
	exists, logs := s.keyExists(key)
	if !exists{
		return "", errors.New("key not found")
	}
	if index > len(logs){
		return "", errors.New("index out of range")
	}
	return logs[index], nil
}

// This is the state machine approach:
// A State machine has state variables encoding its state and
// commands which transform its state.
func (s *Store)HandleCommand(command string, argument []string) (bool, string){
	switch command{
	case "GET":
		value, err := s.Get(argument[0])
		if err != nil{
			return false, ""
		}
		return true, value
	case "SET":
		// Set command for adding new keys
		key := argument[0]
		value := argument[1]
		s.Append(key, value)
		return true, key
	default:
		// Command not known
		return false, ""
	}
}