FROM rockylinux:9.1

LABEL org.opencontainers.image.source https://github.com/bxffour/crest-countries

RUN dnf -y upgrade
# installing psql, tar and migrate
RUN dnf install -y postgresql
RUN dnf install -y tar
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /bin/migrate
   
RUN mkdir /app
RUN mkdir /migrations
WORKDIR /app

COPY ./scripts/db-setup.sh ./entrypoint.sh
COPY ./scripts/wait-for-it.sh /bin/wait-for-it
COPY ./migrations /migrations
COPY ./bin/bs /bin/bootstrap

RUN chmod +x entrypoint.sh

CMD ./entrypoint.sh