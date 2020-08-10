# --- Base ----
FROM golang:1.13-stretch AS base
WORKDIR $GOPATH/src/gitlab.com/betting-game

# ---- Dependencies ----
FROM base AS dependencies
ENV GO111MODULE=on
COPY . .

# ---- Build ----
FROM dependencies AS build
ARG VERSION
ARG BUILD
ARG DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/betting-game -ldflags "-X main.version=${VERSION} -X main.build=${BUILD} -X main.date=${DATE}" ./cmd/server/main.go
# ENTRYPOINT ["/go/bin/betting-game"]

# ---- Base Node ----
FROM node:12-stretch-slim AS nodebase
ENV NODE_ENV=development
RUN mkdir /app
RUN mkdir /app/frontend
WORKDIR /app
RUN npm set progress=false && npm config set depth 0

FROM nodebase AS front-dependencies
WORKDIR /app/frontend
COPY ./frontend/package*.json ./
RUN npm ci && npm cache clean --force

# ---- Front ----
FROM front-dependencies AS front
WORKDIR /app/frontend
COPY ./frontend ./
ENV NODE_ENV=production
RUN npm run build


# ---- Image ----
FROM alpine AS image
COPY --from=front /app/frontend/dist ./frontend/dist
COPY --from=build /go/bin/betting-game /betting-game
ENTRYPOINT ["/betting-game"]
