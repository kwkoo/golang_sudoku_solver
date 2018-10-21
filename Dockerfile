FROM golang:1.10.2 as builder
LABEL builder=true
COPY src /go/src/
RUN set -x && \
	cd /go/src/solver/cmd/solver && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/solver .

FROM scratch
LABEL maintainer="glug71@gmail.com"
LABEL builder=false
COPY --from=builder /go/bin/solver /

# we need to copy the certificates over because we're connecting over SSL
COPY --from=builder /etc/ssl /etc/ssl

# copy timezone info
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /usr/share/zoneinfo/Asia/Singapore /etc/localtime

EXPOSE 8080

ENTRYPOINT ["/solver"]

