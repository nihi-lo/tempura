FROM node:22 AS base

ENV TZ=Asia/Tokyo

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci
