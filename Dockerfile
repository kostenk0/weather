FROM ubuntu:latest
LABEL authors="kostenk0"

ENTRYPOINT ["top", "-b"]