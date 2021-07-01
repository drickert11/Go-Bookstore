FROM postgres:alpine

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD password  
ENV POSTGRES_DB app_staging 

COPY create_database.sql /docker-entrypoint-initdb.d/
COPY create_table.sql /docker-entrypoint-initdb.d/