FROM golang:1.25-alpine AS build-stage

WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux go build -o hazel-bot .

# Create empty JSON files if they don't exist
RUN touch birthdays.json || true
RUN mkdir -p internal/agent && touch internal/agent/agent.json internal/agent/agent_card.json || true

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary
COPY --from=build-stage /app/hazel-bot .

# Create directory structure 
RUN mkdir -p ./internal/agent

# Copy files from build stage - these will exist now due to touch commands
COPY --from=build-stage /app/internal/agent/ ./internal/agent/
COPY --from=build-stage /app/birthdays.json ./
COPY --from=build-stage /app/birthday_workflow.json ./

EXPOSE 3000

CMD [ "./hazel-bot" ]
