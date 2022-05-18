# gowol
wake on lan web written in golang

# TODO
1. add bootstrap theme
2. create api server with go
3. integrate api server with golang library
4. run in docker container
5. docker build . -t gowol
6. push to docker hub
   1. docker image tag gowol:latest developertyrone/gowol:latest
   2. docker push developertyrone/gowol:latest



# QUICKFIX
1. ../../../../../go/pkg/mod/golang.org/x/sys@v0.0.0-20200116001909-b77594299b42/unix/zsyscall_darwin_arm64.go:121:3: //go:linkname must refer to declared function or variable
   1. go get -u golang.org/x/sys