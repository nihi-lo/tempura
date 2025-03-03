FROM node:22

ENV TZ=Asia/Tokyo
ENV NODE_ENV=development

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

EXPOSE 3000
CMD ["npm", "run", "dev"]
