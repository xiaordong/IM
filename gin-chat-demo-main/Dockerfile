FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main .

FROM scratch

COPY --from=builder /app/main /main
COPY --from=builder /app/conf/config.ini /conf/config.ini
COPY --from=builder /app/conf/locales/zh-cn.yaml /conf/locales/zh-cn.yaml

WORKDIR /

EXPOSE 8080

CMD ["/main"]