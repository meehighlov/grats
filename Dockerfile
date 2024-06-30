FROM golang:1.22 as build

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/grats

FROM scratch
COPY  --from=build /bin/grats /bin/grats

CMD ["/bin/grats"]
