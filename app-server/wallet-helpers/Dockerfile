# Installer
FROM node:12 as installer
WORKDIR /usr/src/app
COPY . .
# https://stackoverflow.com/questions/18136746/npm-install-failed-with-cannot-run-in-wd#comment89313677_48170041
RUN npm set unsafe-perm true && npm install

# Builder
FROM node:12 as builder
WORKDIR /usr/src/app
COPY --from=installer /usr/src/app .
RUN npm set unsafe-perm true && npm run build

# App
FROM node:12-alpine as app
LABEL MAINTAINER="Tuan Nguyen <kubeplusplus@gmail.com>"
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/.env .
COPY --from=builder /usr/src/app/lib ./lib
COPY --from=builder /usr/src/app/package.json .
RUN npm install --production
RUN mkdir logs
CMD ["npm", "run", "start"]