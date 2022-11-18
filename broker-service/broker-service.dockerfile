FROM alpine:latest

RUN mkdir /app
COPY cmd/brokerService /app
CMD [ "/app/brokerService" ]