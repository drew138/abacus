FROM migrate/migrate:v4.19.1
COPY *.sql /migrations/
ENTRYPOINT ["migrate", "-path", "/migrations"]
CMD ["--help"]
