package repository

import "sync-bot/storages"

type Repositories struct {
	Command CommandRepository
	User    UserRepository
	Sync    SyncRepository
}

func NewRepositories(db *storages.Database) *Repositories {
	return &Repositories{
		Command: NewCommandRepository(db),
		User:    NewUserRepository(db),
		Sync:    NewSyncRepository(db),
	}
}
