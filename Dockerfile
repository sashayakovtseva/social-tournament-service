FROM alpine

COPY st-service /go/
WORKDIR /go
CMD ["./st-service"]