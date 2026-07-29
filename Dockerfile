FROM golang:1.26.1-alpine as builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags="-s -w" \
    -o /opt/taskmgr ./main.go


FROM scratch

WORKDIR /

COPY --from=builder /opt/taskmgr /taskmgr
COPY --from=builder /src/web ./web

EXPOSE 7540
CMD ["/taskmgr"]