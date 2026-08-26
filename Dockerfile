FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/mpc-coordinator ./cmd/mpc-coordinator
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/mpc-coordinator /mpc-coordinator
EXPOSE 8080
ENTRYPOINT ["/mpc-coordinator"]
