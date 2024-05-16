build: 
	go build -o gofy.exe ./src

secret = 36071D80F000B80B4FB4440BA5F1821B52FEAC7003773C864D4C35D2497CF379

run: 
	go build -o gofy.exe ./src
	@./gofy -jwt-secret=$(secret) -dsn-usr=gofy -dsn-db=gofy -dsn-pwd=freeroam
