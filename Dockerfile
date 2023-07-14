FROM golang:1.21rc2-alpine3.18 as build

WORKDIR /app

RUN apk add --update --no-cache make
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN make build

FROM alpine as main

COPY --from=build /app/app /
RUN apk add --update --no-cache ffmpeg make
RUN wget -P /bin https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp
RUN chmod a+x /bin/yt-dlp

CMD ["make", "run"]
