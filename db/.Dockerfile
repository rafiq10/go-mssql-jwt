ARG DB_PWD=pwd
FROM postgres
ENV POSTGRES_USER bilrafal
ENV POSTGRES_PASSWORD ${DB_PWD}
ENV POSTGRES_DB auth
ENV POSTGRES_HOST_AUTH_METHOD trust
COPY init-user-db.sh /docker-entrypoint-initdb.d/
# ENTRYPOINT [ "docker-entrypoint-initdb.d/init-user-db.sh" ]
# ENV PGDATA=/var/lib/postgresql/data/pgdata 
# VOUME /custom/mount:/var/lib/postgresql/data 
# COPY ./init_db.sql /docker-entrypoint-initdb.d/

# RUN 'ECHO "host all all all " >> pg_hba.conf'
EXPOSE 5432