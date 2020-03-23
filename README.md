# GoLang Mux, Gorm, JWT REST API Boilerplate 

The purpose of this web app is for me to learn Go REST API Development using Mux (router), Gorm (ORM), JSON Web Tokens (JWT), and Golang development in general.

**This is for educational purposes only and probably unsuitable for production** 

## What's included

Basic CRUD routes for user management  

* Show Users `GET /users`
* Show User `GET /users/{userId}`
* Create User `POST /users`
* User Login `POST /users/login`
* Delete User `DELETE /users/{userId}`
* Update User `PUT /users/{userId}` * Note only the user can update their own name

Several routes are protected and require JWT tokens, which can be generated using the login route.
You will need to create a user by sending a post request to the createUser route.

## Configuration

Make sure to copy `.env.sample` to `.env` and update all fields (DB fields are important!)

**Please note that this is using the MySQL driver, if you prefer to use another driver, update `db.go` accordingly**

Gorm is setup to automigrate the Users table, so it should be plug and play.

## Installation

Make sure to have all required external deps. Look at Godeps config file to view them all.

**Run:**

To run using Go: `go run main.go`

To view application in browser: `localhost:3001 or locahost:YOUR_PORT_ENV (go run)`
  
## Todos
 
[] Add dependency management<br>
[] Implement redis (not know why but just because)<br>
[] Gorm setupo for migration<br>
[] Add Authorization with Token<br>
