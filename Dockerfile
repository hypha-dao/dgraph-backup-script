FROM golang:1.16
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
RUN mkdir /app
WORKDIR /app
COPY . /app
RUN go install
