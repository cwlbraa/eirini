FROM golang:1.13 as builder
WORKDIR /eirini/
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux go build -mod vendor -trimpath -a -installsuffix cgo -o staging-reporter ./cmd/staging-reporter
ARG GIT_SHA
RUN if [ -z "$GIT_SHA" ]; then echo "GIT_SHA not set"; exit 1; else : ; fi

FROM scratch
COPY --from=builder /eirini/staging-reporter /usr/local/bin/staging-reporter
USER 1001
ENTRYPOINT [ "/usr/local/bin/staging-reporter", \
	"--config", \
	"/etc/eirini/config/staging-reporter.yml" \
]
ARG GIT_SHA
LABEL org.opencontainers.image.revision=$GIT_SHA \
      org.opencontainers.image.source=https://code.cloudfoundry.org/eirini
