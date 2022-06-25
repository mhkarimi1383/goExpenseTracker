FROM golang:1.18.2-alpine3.16 as builder

RUN apk add --no-cache git

COPY . /go/src/github.com/mhkarimi1383/goExpenseTracker
WORKDIR /go/src/github.com/mhkarimi1383/goExpenseTracker

## we have vendor directory in our project no need to get packages again
# RUN go get -v ./...

RUN go build -o /goExpenseTracker

FROM alpine:3.14 as runtime

WORKDIR /app

## copy and prepare binary file
COPY --from=builder /goExpenseTracker .
RUN chmod +x ./goExpenseTracker

## copy static files
COPY templates ./templates
COPY translate.yaml ./translate.yaml
## making it non-root user
RUN adduser -D no-name
USER no-name:no-name

CMD ["/app/goExpenseTracker"]