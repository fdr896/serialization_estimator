FROM golang:1.19

WORKDIR /usr/local/go/src/serialization_estimator

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/spf13/cobra-cli@latest

COPY . .

RUN cd bin/estimator && go install

ARG PORT
ARG METHOD
ENV PORT ${PORT}
ENV METHOD ${METHOD}

CMD estimator
