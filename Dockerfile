FROM golang:1.23 AS modules

WORKDIR /app

# Download Go modules
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

FROM modules AS build

ENV CGO_ENABLED=0

ENV GOOS=linux

# Copy src code
COPY main.go ./

# Build
RUN go build -o /sfs

FROM gcr.io/distroless/base-debian12 AS runtime

COPY --from=build /sfs /sfs

EXPOSE 3000

USER nonroot:nonroot

# Run
CMD [ "/sfs" ]