# ##############################
# # Stage 1: Builder           #
# ##############################
# FROM golang:alpine AS builder
# WORKDIR /app

# # Copy Go modules files
# COPY go.mod go.sum ./

# # Download dependencies
# RUN go mod download

# # Copy the rest of the application source code
# COPY . .

# # Build the application
# RUN go build -o main .

# ##############################
# # Stage 1: Builder           #
# ##############################
# FROM scratch
# COPY --from=builder /app/main /main

# EXPOSE 8000

# ENTRYPOINT ["/main"]

FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]