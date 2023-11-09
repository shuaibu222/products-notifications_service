FROM alpine:latest

RUN mkdir /app

COPY notificationsApp /app

CMD [ "/app/notificationsApp"]