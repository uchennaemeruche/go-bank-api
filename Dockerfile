
# # 1.18.4-alpine3.16
FROM golang:1.18-alpine3.16 as buildStage
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.16
WORKDIR /app
COPY --from=buildStage /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]