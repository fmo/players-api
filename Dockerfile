FROM alpine:latest

RUN mkdir /app

COPY playersApi /app

CMD [ "/app/playersApi" ]
