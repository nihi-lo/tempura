FROM node:22 AS base

ENV NODE_ENV=development
ENV TZ=Asia/Tokyo

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

EXPOSE 5173
CMD ["npm", "run", "dev"]
