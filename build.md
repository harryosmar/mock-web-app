```shell script
env GOOS=darwin GOARCH=amd64 go build -o build/mock-web-app-darwin-64 .\
&& env GOOS=windows GOARCH=amd64 go build -o build/mock-web-app-windows-64 .\
&& env GOOS=windows GOARCH=386 go build -o build/mock-web-app-windows-386 . \
&& chmod +X build/
```