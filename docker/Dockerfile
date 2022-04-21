FROM golang:1.18 as build
COPY . /folduka
WORKDIR /folduka
RUN   GO111MODULE=on go mod tidy -v && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o folduka main.go actions.go login.go display.go workflow.go download.go sign.go signage.go smtp.go utils.go websocket.go template.go database.go queue.go

FROM alpine:3.12 as downloads
ADD https://github.com/hairyhenderson/gomplate/releases/download/v3.8.0/gomplate_linux-amd64-slim /gomplate
ADD https://github.com/reconquest/orgalorg/releases/download/1.0/orgalorg_1.0_linux_amd64.tar.gz /orgalorg.tar.gz
ADD https://github.com/xo/usql/releases/download/v0.7.8/usql-0.7.8-linux-amd64.tar.bz2 /usql.tar.bz2
ADD https://get.helm.sh/helm-v3.3.1-linux-amd64.tar.gz /helm.tar.gz
RUN apk add bzip2 curl
RUN tar jxvf /usql.tar.bz2 && tar zxvf /orgalorg.tar.gz && tar zxvf /helm.tar.gz
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && chmod +x ./kubectl 

FROM ubuntu:20.04
RUN apt-get update && apt-get install -y software-properties-common && add-apt-repository -y ppa:alex-p/tesseract-ocr
RUN apt-get update && apt-get install -y jq ca-certificates curl ffmpeg git openjdk-8-jre-headless tesseract-ocr-all netcat
RUN update-ca-certificates
RUN ssh-keygen -t rsa -N "" -f ~/.ssh/id_rsa
#
COPY --from=build /folduka/folduka /folduka
COPY --from=downloads /gomplate /usr/local/bin/gomplate
COPY --from=downloads /orgalorg /usr/local/bin/orgalorg
COPY --from=downloads /usql /usr/local/bin/usql
COPY --from=downloads /linux-amd64/helm /usr/local/bin/helm
COPY --from=downloads /kubectl /usr/local/bin/kubectl
#
ADD templates /templates
ADD public /public
ADD signpdf/lib /signpdf/lib
ADD signpdf/PortableSigner.jar /signpdf/PortableSigner.jar
CMD [ "/folduka" ]
