FROM debian:stable-slim
RUN apt-get update -y
RUN apt-get install -y --no-install-recommends x11-apps xauth
RUN apt-get clean
ENTRYPOINT ["xeyes"]
