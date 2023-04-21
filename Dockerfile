FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.20.3 AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /ddns

COPY go.mod /ddns
COPY go.sum /ddns
COPY Makefile /ddns
COPY LICENSE /ddns
COPY README.md /ddns
COPY VERSION /ddns
COPY **.go /ddns
RUN export VERSION=$(cat VERSION) && \
    CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-X 'main.Version=v${VERSION}'" -a -installsuffix cgo -o ddns .


FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot
WORKDIR /ddns

COPY --chown=nonroot --from=builder /ddns/ddns /ddns
COPY --chown=nonroot --from=builder /ddns/README.md /ddns

CMD ["/ddns/ddns"]

USER nonroot:nonroot
