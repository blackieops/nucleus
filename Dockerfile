FROM golang:1.17.7
ENV CGO_ENABLED 0
ADD . /src
WORKDIR /src
RUN go build -a --installsuffix cgo --ldflags="-s" -o nucleus

FROM debian:11-slim
RUN apt-get update && \
	apt-get upgrade && \
	apt-get install -y ca-certificates && \
	apt-get clean
ENV GIN_MODE=release
COPY --from=0 /src/nucleus /usr/bin/nucleus
ENTRYPOINT ["/usr/bin/nucleus"]
CMD ["-config", "/etc/nucleus/config.yaml"]
