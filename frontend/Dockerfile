FROM node:18-alpine as build
WORKDIR /app

COPY package*.json ./
RUN npm install
COPY . .

# ✅ Accept environment variables from CapRover build args
ARG REACT_APP_BACKEND_GO
ARG REACT_APP_BACKEND_PYTHON

# ✅ Expose them as ENV so React can use them at build time
ENV REACT_APP_BACKEND_GO=$REACT_APP_BACKEND_GO
ENV REACT_APP_BACKEND_PYTHON=$REACT_APP_BACKEND_PYTHON

# ✅ Build React app with the injected environment variables
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
