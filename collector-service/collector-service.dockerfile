FROM alpine:latest

RUN mkdir /app
COPY cmd/collectorService /app
CMD [ "/app/collectorService" ]