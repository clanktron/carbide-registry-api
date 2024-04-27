FROM golang:1.22.0-alpine AS builder
COPY . /build
WORKDIR /build
RUN apk add make git 
RUN apk add --no-cache ca-certificates
RUN make dependencies
RUN make
WORKDIR /permissions
RUN echo "carbide-registry-api:x:1001:1001::/:" > passwd && echo "carbide-registry-api:x:2000:carbide-registry-api" > group

FROM scratch
COPY --from=builder /permissions/passwd /etc/passwd
COPY --from=builder /permissions/group /etc/group
USER carbide-registry-api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/carbide-registry-api /
ENTRYPOINT ["/carbide-registry-api"]
