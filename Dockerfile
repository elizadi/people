FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o people main.go
EXPOSE 8000
CMD ["./people"]