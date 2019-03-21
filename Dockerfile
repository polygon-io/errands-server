# Step #1
FROM golang as firststage
WORKDIR /work
ADD . .
# RUN go test ./...
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o errands-server .



# Step #2:
FROM alpine:latest  
ENV TZ=America/New_York
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=firststage /work/errands-server .
CMD ["./errands-server"]