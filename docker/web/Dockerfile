FROM node:23-alpine3.20 AS build-env

WORKDIR /usr/src/app
COPY . ./

RUN npm install
RUN npm run build

EXPOSE 8080 

CMD [ "npm", "run", "start" ]