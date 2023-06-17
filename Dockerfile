FROM golang:1.20.5

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

CMD ["make", "run"]
