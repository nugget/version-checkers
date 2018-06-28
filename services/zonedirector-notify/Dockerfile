FROM golang:1.10.1 AS builder

ARG dockerhost
ARG dockeruser
ARG dockerver
ENV dockerhost=$dockerhost
ENV dockeruser=$dockeruser
ENV dockerver=$dockerver

# # Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/nugget/version-checkers
COPY . .

WORKDIR $GOPATH/src/github.com/nugget/version-checkers/services/zonedirector-notify
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /zonedirector-notify .

FROM scratch
COPY --from=builder /zonedirector-notify ./
ADD ./extras/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./zonedirector-notify"]

