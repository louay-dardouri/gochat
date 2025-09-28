FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./server -ldflags "-w -s" ./cmd/server/

FROM alpine

RUN apk add --no-cache openssh busybox-extras shadow

COPY --from=builder /app/server /usr/local/bin/chat-server

RUN adduser -D username && echo "username:username" | chpasswd

COPY guest_profile.sh /home/username/.profile
RUN chown username:username /home/username/.profile && chmod 755 /home/username/.profile

RUN mkdir /var/run/sshd

RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config
RUN sed -i 's/#PermitEmptyPasswords no/PermitEmptyPasswords no/' /etc/ssh/sshd_config

COPY server_entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 22 8080

ENTRYPOINT [ "/entrypoint.sh" ]
