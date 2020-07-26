# What is this project?

## Backend
* Dummy Shop with in-memory storage. You can replace the implementation
behind the interfaces wit "real" databases.
* Uses golang.
* Does not have a Dockerfile, because we do not need that for this demo
* Empty kubernetes directory, here we would put helm charts to describe/deploy this application
* Some unit tests
* No integration tests

## Frontend
I usually dont do "frontend".

A demo show-casing my VueJS skills that I acquired hacking this project.
Don't expect much. 
I intentionally did not use a whole vuejs-cli webpack pipeline, because
it was simply to much bloat for a demo project.

* Uses VueJS
* Plain JS
* Single file index.js

# How to use?

1. `cd backend/src/cmd/server`
2. `go run main.go`
3. `open http://localhost:8181/`
4. Signup
5. Login

