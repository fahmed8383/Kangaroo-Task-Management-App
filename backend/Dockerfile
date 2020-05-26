# use base alpine golang image
FROM golang:alpine AS builder

# create a /app directory and make that the working directory for the builder container
RUN mkdir /app
WORKDIR /app

# add git for go vendor imports to the container
RUN apk add --no-cache git

# copy the contents of the current directory to the builder container
COPY . .

# arguement from the docker compose file used to differenciate between the different services
ARG service

# move the work directory to the current service we want to containerize
WORKDIR /app/${service}

# build the go app and save it as /service exe
RUN CGO_ENABLED=0 GOOS=linux go build -o /service

# pull alpine image for production
FROM alpine:latest AS production

# copy the /service exe from the builder to the production file
COPY --from=builder /service /service

# run the go app exe 
CMD ["/service"]