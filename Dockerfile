FROM alpine:3.11 as app

RUN apk add --no-cache libc6-compat

WORKDIR /app
COPY api.yaml /app/api.yaml
COPY build/blog-aggregator /app/blog-aggregator
COPY LICENSE /app/LICENSE
COPY NOTICE /app/NOTICE
RUN chmod +x /app/blog-aggregator

EXPOSE 80
ENTRYPOINT [ "/app/blog-aggregator" ]
CMD ["-port=80"]