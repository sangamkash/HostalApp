FROM golang:1.24.3-alpine AS build

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Generate docs (simplified for Fiber)
RUN swag init --generalInfo ./cmd/api/main.go --output ./docs

# Verify docs
RUN ls -la ./docs && [ -f ./docs/swagger.json ]

# Build the application
RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/docs /app/docs

# Make sure your main application is configured to serve from ./docs
ENV SWAGGER_JSON=/app/docs/swagger.json
EXPOSE ${PORT}
CMD ["./main"]

# Frontend builds remain the same
FROM node:20 AS frontend_builder
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/. .
RUN npm run build

FROM node:23-slim AS frontend
RUN npm install -g serve
COPY --from=frontend_builder /frontend/dist /app/dist
EXPOSE 5173
CMD ["serve", "-s", "/app/dist", "-l", "5173"]