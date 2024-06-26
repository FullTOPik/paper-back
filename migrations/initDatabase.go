package main

import (
	"paper_back/config"
)

func main() {
	config.Connect()
	defer config.Disconnect()

	if _, err := config.Database.Exec(`
		CREATE TABLE users (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			username VARCHAR(255) NOT NULL,
			role ENUM('USER', 'ADMIN') NOT NULL,
			created_at Date DEFAULT NOW(),
			last_visit Date,
			password VARCHAR(255) NOT NULL,
			CONSTRAINT UC_User UNIQUE (username)
		);
		`); err != nil {
		panic(err)
	}

	if _, err := config.Database.Exec(`
	CREATE TABLE user_codes (
		code VARCHAR(255) NOT NULL unique,
		user_id INT NOT NULL unique,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);
	`); err != nil {
	panic(err)
}

	if _, err := config.Database.Exec(`
		CREATE TABLE tokens (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			user_id INT NOT NULL, 
			access VARCHAR(255) NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);
		`); err != nil {
		panic(err)
	}

	if _, err := config.Database.Exec(`
		CREATE TABLE contacts (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			initiator_id INT NOT NULL,
			agreeing_id INT NOT NULL,
			secret VARCHAR(255) NOT NULL unique,
			status ENUM('ACTIVE', 'DELETED', 'CREATED') NOT NULL,
			created_at Date DEFAULT NOW(), 
			updated_at Date,
			CONSTRAINT `fk_contact_initiator` FOREIGN KEY (initiator_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE RESTRICT,
			CONSTRAINT `fk_contact_agreeing` FOREIGN KEY (agreeing_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE RESTRICT
		);
		`); err != nil {
		panic(err)
	}

	if _, err := config.Database.Exec(`
		CREATE TABLE messages (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			contact_id INT NOT NULL,
			sender ENUM('INITIATOR', 'AGREEING') NOT NULL,
			text TEXT NOT NULL,
			created_at Date DEFAULT NOW(),
			updated_at Date,
			CONSTRAINT `fk_message_contact` FOREIGN KEY (contact_id) REFERENCES contacts (id) ON DELETE CASCADE ON UPDATE RESTRICT
		);
		`); err != nil {
		panic(err)
	}

	if _, err := config.Database.Exec(`
		CREATE TABLE groups (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			owner_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at Date DEFAULT NOW(),
			CONSTRAINT UC_Group UNIQUE (name),
			FOREIGN KEY (owner_id) REFERENCES users (id)
		);
		`); err != nil {
		panic(err)
	}

	if _, err := config.Database.Exec(`
		CREATE TABLE posts (
			id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
			title VARCHAR(255) NOT NULL,
			group_id INT NOT NULL,
			text TEXT NOT NULL,
			premium BOOL DEFAULT false,
			created_at Date DEFAULT NOW(),
			updated_at Date
		);
		`); err != nil {
		panic(err)
	}
}
