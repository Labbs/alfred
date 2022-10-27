FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY app /app
COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh && chmod +x /app
ENTRYPOINT ["/entrypoint.sh"]