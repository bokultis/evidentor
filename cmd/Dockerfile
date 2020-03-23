FROM golang:latest
WORKDIR /app
COPY bin/evidentor /app
ADD auth_model.conf /app
ADD .env /app
ADD info.log /app
CMD ["./evidentor"]
EXPOSE 8080