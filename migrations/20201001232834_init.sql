-- +goose Up
CREATE TABLE IF NOT EXISTS users(name TEXT, rank INTEGER, PRIMARY KEY(name));
CREATE TABLE IF NOT EXISTS chat_history(timestamp INTEGER, username TEXT, msg TEXT);
CREATE TABLE IF NOT EXISTS command_rank(command VARCHAR(10), rank INTEGER);

