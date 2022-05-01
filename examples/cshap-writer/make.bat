set GOARCH=amd64
set CGO_ENABLED=1
go build -ldflags "-s -w" -x -buildmode=c-shared -o Spring.MMDBWriter.dll main_cshap.go