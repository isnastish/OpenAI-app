FROM golang:1.23 AS build-env

WORKDIR /openai/service/

ADD . /openai/service/

RUN CGO_ENABLED=0 GOOS=linux go build -a -v -o /openai/bin/service github.com/isnastish/openai/service/

FROM golang:1.23-alpine3.21 AS run-env

COPY --from=build-env /openai/bin/service/ /openai/service/ 

EXPOSE 3030 

CMD [ "/openai/service/service" ]  

