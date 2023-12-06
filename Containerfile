FROM docker.io/golang:1.21-alpine
LABEL maintainer="Mateus Melchiades"

WORKDIR /home/user

# Install dependencies
RUN apk add --update gcc musl-dev sqlite make

# Copy project files and build
COPY main.go go.mod go.sum go.work Makefile /home/user/
COPY core /home/user/core/
COPY diff /home/user/diff/
COPY types /home/user/types/
COPY vendor /home/user/vendor/
RUN make

# Cleanup
RUN rm -rf main.go go.mod go.sum go.work Makefile core/ types/ vendor/

EXPOSE 8080

CMD ["/home/user/differ"]
