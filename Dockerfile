FROM postgres:10.3

ARG organizations=""

COPY schema/*.sql /docker-entrypoint-initdb.d/

RUN for org in $organizations; do \
  sed "s/organization/$org/g" /docker-entrypoint-initdb.d/organization.sql > /docker-entrypoint-initdb.d/"$org".sql; \
  done

EXPOSE 5432/tcp

VOLUME /var/lib/postgresql/data
