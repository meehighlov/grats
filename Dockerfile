FROM golang:1.22 as build

WORKDIR /app

COPY . /app

RUN go mod download

COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o grats

FROM scratch
COPY  --from=build /app/grats /grats

CMD ["/app/grats"]
