FROM mysql:latest
LABEL maintainer="bob@2kranki.us"
ENV MYSQL_ROOT_PASSWORD="Passw0rd!"
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY ./app .
WORKDIR /
#CMD "/usr/src/app/entrypoint.sh"
