FROM alpine

COPY ./reaper /usr/local/bin/

ENTRYPOINT reaper
