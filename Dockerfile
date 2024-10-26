FROM python:3.13-alpine

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

WORKDIR /app

RUN apk add --no-cache supervisor nginx uwsgi uwsgi-python3 npm

COPY app.py package.json schema.sql tailwind.config.js requirements.txt /app/
COPY htm /app/htm/
COPY py /app/py/
COPY static /app/static/
COPY docker /app/

COPY docker/nginx.conf /etc/nginx/nginx.conf

RUN pip install --no-cache-dir -r requirements.txt

ENV SM_DATA=/data

RUN adduser nginx nginx && mkdir /data && chown nginx:nginx /data && npm install && npx tailwindcss -i ./static/css/input.css -o ./static/css/main.css


EXPOSE 80/tcp
VOLUME /data
VOLUME /corsa

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/app/supervisord.conf"]
