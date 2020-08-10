# ---- Backend Base ----
FROM golang:1.13-stretch AS base
WORKDIR $GOPATH/src/gitlab.com/betting-game

# ---- Backend Dependencies ----
FROM base AS dependencies
ENV GO111MODULE=on
COPY . .

# ---- Backend test dependencies ----
FROM dependencies AS test-dependencies
ENV GO111MODULE=on
RUN go get -u github.com/axw/gocov/gocov && GO111MODULE=off go get -u gopkg.in/matm/v1/gocov-html

# ---- Backend Test ----
FROM test-dependencies AS test
RUN go test -v -cpu 1 -failfast -coverprofile=coverage.out -covermode=set ./...
RUN gocov convert coverage.out | gocov-html > /index.html
RUN grep -v "_mock" coverage.out >> filtered_coverage.out
RUN go tool cover -func filtered_coverage.out

# ---- Backend test dependencies ----
FROM dependencies AS lint-dependencies
ENV GO111MODULE=on
RUN go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

# ---- Lint ----
FROM lint-dependencies AS lint
RUN golangci-lint run -c ./.golangci.yml

# ---- Backend Build ----
FROM dependencies AS build
ARG VERSION
ARG BUILD
ARG DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/betting-game -ldflags "-X main.version=${VERSION} -X main.build=${BUILD} -X main.date=${DATE}" ./cmd/server/main.go

# ---- Frontend Base ----
FROM node:12-stretch-slim AS nodebase
ENV NODE_ENV=development
RUN mkdir /app
RUN mkdir /app/frontend
WORKDIR /app
RUN npm set progress=false && npm config set depth 0

# ---- Frontend Dependencies ----
FROM nodebase AS front-dependencies
WORKDIR /app/frontend
COPY ./frontend/package*.json ./
RUN npm ci && npm cache clean --force

# ---- Frontend Build ----
FROM front-dependencies AS front
WORKDIR /app/frontend
COPY ./frontend ./
ENV NODE_ENV=production
RUN npm run build


# ---- Final Image - Both builds ----
FROM alpine AS image
COPY --from=front /app/frontend/dist ./frontend/dist
COPY --from=build /go/bin/betting-game /betting-game
ENTRYPOINT ["/betting-game"]

