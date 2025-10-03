# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built application
COPY --from=builder /app/build /usr/share/nginx/html

# Copy nginx configuration template
COPY nginx.conf /etc/nginx/conf.d/default.conf.template

# Create startup script to handle PORT env var
RUN echo '#!/bin/sh' > /start.sh && \
    echo 'PORT=${PORT:-80}' >> /start.sh && \
    echo 'sed "s/listen 80;/listen $PORT;/g" /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf' >> /start.sh && \
    echo 'nginx -g "daemon off;"' >> /start.sh && \
    chmod +x /start.sh

ENTRYPOINT ["/start.sh"]
