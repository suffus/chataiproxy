FROM golang:1.22.1-alpine  as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /chat-ai

EXPOSE 8080
ENTRYPOINT ["/chat-ai"]
CMD [ "-openAIKey" ${OPENAI_KEY} ]