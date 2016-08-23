# go-redis-delete-pattern
Delete keys based on pattern passed as argument
->Requirements: Have go installed, with gopath setup.
* `cp config.sample config.go` // fill redis credentials in config.go
* `go build -o redisdeletepattern` //this will build the project
* `./redisdeletepattern "patternyouwanttodelete"` //this will output the total keys deleted, error if occured.
