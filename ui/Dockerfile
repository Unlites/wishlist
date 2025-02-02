FROM node:23 AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

RUN npm run build

FROM nginx:1.23
COPY --from=builder /app/dist /usr/share/nginx/html
