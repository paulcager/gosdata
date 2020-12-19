FROM paulcager/go-base:latest as build
WORKDIR /go/src/

COPY . /go/src/github.com/paulcager/gosdata
RUN cd /go/src/github.com/paulcager/gosdata && go test ./... && go install ./...

####################################################################################################


FROM debian:stable-slim
RUN apt-get update && apt-get -y upgrade && apt-get install -y ca-certificates
WORKDIR /app
COPY --from=build /go/bin/* .
RUN ./load-gosdata -o . && rm terrain-50.zip
EXPOSE 9091
CMD ["/app/serve-gosdata", "--port", ":9091", "-d", "data" ]

