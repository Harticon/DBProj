# Use an official Python runtime as a parent image

FROM golang:latest AS builder


LABEL maintainer="Vojtěch Hromádka<hromadkavojta@gmail.com>"


#RUN mkdir -p /go/src/app
#WORKDIR /go/src/app

WORKDIR /go/src/github.com/Harticon/DBproj

# Copy everything from present working dir
COPY . .

#Golang dependencies go get/go mod init

#RUN GO111MODULE=on go mod vendor

#WORKDIR /go/src/app/user/cmd/user

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./cmd

#RUN ls
# Run app.py when the container launches
#ENTRYPOINT ["./grpcProj"]

# Make port 80 available to the world outside this container
#EXPOSE 8080

FROM alpine:latest

RUN ls
WORKDIR /go/bin
ENV PATH=/bin

RUN ls

COPY --from=builder /go/src/github.com/Harticon/DBproj/cmd/prod.db .
COPY --from=builder /go/bin/cmd .

CMD ["./cmd"]