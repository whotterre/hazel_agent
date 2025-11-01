FROM golang:1.25-alpine AS build-stage

WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux go build -o hazel-bot .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=build-stage /app/hazel-bot .

EXPOSE 3000

CMD [ "./hazel-bot" ]
