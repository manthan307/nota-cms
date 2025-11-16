# 1) Build stage
FROM node:22.19.0-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

# Build static export
RUN npm run build

EXPOSE 3000

CMD ["npm","start"]
