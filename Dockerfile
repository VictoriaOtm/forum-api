FROM ubuntu:16.04

LABEL author="Victoria Kamoldinova"

# Actually pgver is 10.3
ENV PGVER 10
ENV GOVER 1.10.1

# Installing postgres
RUN apt-get -y update
RUN apt-get install -y wget git
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ xenial-pgdg main" >>  /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y update
RUN apt-get install -y postgresql-$PGVER wget git

USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker && /etc/init.d/postgresql stop
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "synchronous_commit='off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Installing golang
USER root
RUN wget https://dl.google.com/go/go$GOVER.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go$GOVER.linux-amd64.tar.gz && mkdir go go/src go/bin go/pkg

ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

## Installing forum api
COPY ./ $GOPATH/src/github.com/VictoriaOtm/forum-api/
WORKDIR $GOPATH/src/github.com/VictoriaOtm/forum-api/
RUN go install github.com/VictoriaOtm/forum-api

# Exposing ports
EXPOSE 5432

# running
USER postgres
CMD /etc/init.d/postgresql start && forum-api
