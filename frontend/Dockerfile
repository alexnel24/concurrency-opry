FROM node:20 AS dev

WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .

EXPOSE 5173

FROM dev AS build

RUN npm run build

FROM nginx:alpine AS prod

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
