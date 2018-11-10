FROM ubuntu:18.04

# Обвновление списка пакетов
ARG USERNAME
ARG PASSWORD
RUN apt-get -y update
ENV USERNAME $USERNAME
ENV PASSWORD $PASSWORD

#
# Установка postgresql
#
ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf


# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432
USER root
# Add VOLUMEs to allow backup of config, logs and databases
#VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Установка golang
ENV GOVER 1.10
RUN apt-get install -y golang-$GOVER
RUN apt-get install -y git

# Выставляем переменную окружения для сборки проекта
ENV GOROOT /usr/lib/go-$GOVER
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

# Копируем исходный код в Docker-контейнер
WORKDIR $GOPATH/src/techpark-db
COPY . $GOPATH/src/techpark-db

RUN go get -u github.com/gorilla/mux
RUN go get github.com/lib/pq
RUN go install .
EXPOSE 5000
CMD /etc/init.d/postgresql start && techpark-db
