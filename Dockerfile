FROM golang
WORKDIR /go/src/github.com/digitorus/md2csv/
RUN go get -d -v github.com/gomarkdown/markdown
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o md2csv .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/digitorus/md2csv/md2csv .
CMD ["./md2csv"]  