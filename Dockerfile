FROM postgres:10.3

COPY schema/*.sql /docker-entrypoint-initdb.d/

EXPOSE 5432/tcp

VOLUME /var/lib/postgresql/data
