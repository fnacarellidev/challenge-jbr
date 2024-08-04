FROM node:slim as builder

WORKDIR /app

COPY public public
COPY src src
COPY package.json package-lock.json .

RUN npm install
RUN npm run build

FROM node:slim

WORKDIR /app

RUN npm install -g serve
COPY --from=builder /app/build build

ENTRYPOINT ["serve", "-s", "build"]
