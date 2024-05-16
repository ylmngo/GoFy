build: 
	go build -o gofy.exe ./src

run: 
	go build -o gofy.exe ./src
	@./gofy  -dsn-usr=gofy -dsn-db=gofy -dsn-pwd=freeroam
