FROM golang:1.24.4 as builder

ENV BIN_FILE /opt/system-monitor/system-monitor
ENV CODE_DIR /go/src/system-monitor

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
    -ldflags "$LDFLAGS" \
    -o ${BIN_FILE} cmd/monitor/*

FROM debian:stable-slim

LABEL ORGANIZATION="DimRyb"
LABEL SERVICE="system-monitor"
LABEL MAINTAINERS="dimryb@bk.ru"

ENV BIN_FILE /opt/system-monitor/system-monitor
ENV CONFIG_FILE /etc/system-monitor/config.yaml

RUN apt update && \
    apt install -y --no-install-recommends procps sysstat && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder ${BIN_FILE} ${BIN_FILE}
COPY ./configs/monitor.yaml ${CONFIG_FILE}

CMD ["sh", "-c", "${BIN_FILE} -config ${CONFIG_FILE}"]