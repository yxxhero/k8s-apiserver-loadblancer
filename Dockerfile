FROM  golang:1.22-alpine as builder

RUN apk add --no-cache make git
WORKDIR /workspace/k8s-apiserver-loadblancer

COPY go.mod go.sum /workspace/k8s-apiserver-loadblancer/
RUN go mod download

COPY . /workspace/k8s-apiserver-loadblancer
RUN go build -o dist/k8s-apiserver-loadblancer main.go 

# -----------------------------------------------------------------------------

FROM alpine:3.19


COPY --from=builder /workspace/k8s-apiserver-loadblancer/dist/k8s-apiserver-loadblancer /usr/local/bin/k8s-apiserver-loadblancer

CMD ["/usr/local/bin/k8s-apiserver-loadblancer"]
