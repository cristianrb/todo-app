FROM alpine:latest

RUN mkdir /app
COPY cmd/recordService /app
CMD [ "/app/recordService" ]