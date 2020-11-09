# Phoosball

Phoosball score tracking webapp written in Go

## Development Set Up
1. install go
2. start mariadb docker container (`cd db/ && ./run.sh`)
3. install dependencies (`cd ../server && go get -u ./...`)
4. run the server in debug (test) mode (`go run phoos_server.go -debug`)
5. run the react development server (`cd ../ui && npm start`)
