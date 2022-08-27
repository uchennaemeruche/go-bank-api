
# # 1.18.4-alpine3.16
FROM golang:1.18-alpine3.16 as buildStage
WORKDIR /app
COPY . .
RUN go build -o main main.go
# Install Curl and use it to download and extract migrate package.
RUN apk add curl  
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.16
WORKDIR /app
COPY --from=buildStage /app/main .
COPY  --from=builderStage /app/migrate .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./migrations

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]