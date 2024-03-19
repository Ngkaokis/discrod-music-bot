FROM golang:1.20.5

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg libopus-dev libopusfile-dev
RUN wget -P /bin https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp
RUN chmod a+x /bin/yt-dlp

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

CMD ["make", "dev"]
