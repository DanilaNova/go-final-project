FROM ubuntu:latest

WORKDIR /app

COPY --parents web main ./

ENV TODO_PORT=7540
ENV TODO_DBFILE=./scheduler.db
ENV TODO_PASSWORD=test1234

EXPOSE 7540

CMD ["./main"]