# First stage, build the application
FROM node:lts-alpine as build-stage
WORKDIR /app
COPY *.json ./
COPY src ./src
COPY public ./public

RUN npm ci

RUN npm run build

# Official Nginx image (Yes, in this article I always use the `latest`. Kill me!)
FROM nginx:stable-alpine

# Copy build from build-stage
COPY --from=build-stage /app/build /usr/share/nginx/html

# Copy Nginx configuration file to the container
COPY ./nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Command to run
CMD ["nginx", "-g", "daemon off;"]