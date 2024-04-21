FROM golang:1.22.2-alpine as builder

WORKDIR /build

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o attendance

########################################################
### Runner
########################################################
FROM alpine:3.19 AS runner

ARG USER=server
ARG GROUP=$USER
ARG UID=1001
ARG GID=1001
ARG HOME=/home/$USER

ENV CONFIG_FILE=$HOME/config.yml
ENV MIGRATION_DIR=$HOME/migrations

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Makassar

COPY --from=builder /build/attendance /usr/local/bin/attendance

RUN addgroup --gid $GID $GROUP && \
    adduser --disabled-password --ingroup $USER --no-create-home --uid $UID \
    $USER

USER $USER

WORKDIR $HOME

COPY ./store/migrations/ ./store/migrations/

CMD ["sh", "-c", "attendance --config $CONFIG_FILE"]