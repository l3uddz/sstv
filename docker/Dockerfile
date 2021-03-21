FROM sc4h/alpine-s6overlay:3.12

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

ENV \
  PATH="/app/sstv:${PATH}" \
  APP_CONFIG="/config/config.yml" \
  APP_LOG="/config/activity.log" \
  APP_VERBOSITY="0"

# Binary
COPY ["dist/sstv_${TARGETOS}_${TARGETARCH}${TARGETVARIANT:+_7}/sstv", "/app/sstv/sstv"]

# Add root files
COPY ["docker/run", "/etc/services.d/sstv/run"]

# Volume
VOLUME ["/config"]

# Port
EXPOSE 1411