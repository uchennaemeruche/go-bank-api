
# # 1.18.4-alpine3.16
FROM golang:1.18-alpine3.16
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 8080
CMD [ "/app/main" ]