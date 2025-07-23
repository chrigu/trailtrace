# Build Stage
FROM node:20 AS builder
WORKDIR /app
COPY . .
RUN npm install
RUN npm run generate

# Production Stage
FROM nginx:alpine
COPY --from=builder /app/.output/public /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
