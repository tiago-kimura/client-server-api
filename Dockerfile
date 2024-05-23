FROM nouchka/sqlite3
COPY init_db.sh /docker-entrypoint-initdb.d/init_db.sh
RUN chmod +x /docker-entrypoint-initdb.d/init_db.sh
CMD ["/docker-entrypoint-initdb.d/init_db.sh", "&&", "tail", "-f", "/dev/null"]
