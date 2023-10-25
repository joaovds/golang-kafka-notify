package models

import "sync"

type Notification struct {
  From       User      `json:"from"`
  To         User      `json:"to"`
  Message    string    `json:"message"`
}

type UserNotifications map[string][]Notification

type NotificationStore struct {
    data UserNotifications
    mu   sync.RWMutex
}

func NewNotificationStore() *NotificationStore {
  return &NotificationStore{
    data: make(UserNotifications),
  }
}

func (ns *NotificationStore) Add(userID string, notification Notification) {
  ns.mu.Lock()
  defer ns.mu.Unlock()
  ns.data[userID] = append(ns.data[userID], notification)
}

func (ns *NotificationStore) Get(userID string) []Notification {
  ns.mu.RLock()
  defer ns.mu.RUnlock()
  return ns.data[userID]
}

