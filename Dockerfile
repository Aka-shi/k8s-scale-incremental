FROM golang:1.15-alpine

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN go build -o k8s-scaler main.go

ENTRYPOINT [ "/app/k8s-scaler" ]