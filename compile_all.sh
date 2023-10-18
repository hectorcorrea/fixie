GOOS=darwin go build -o fixie
GOOS=linux go build -o fixie_linux
GOOS=windows GOARCH=386 go build -o fixie.exe
