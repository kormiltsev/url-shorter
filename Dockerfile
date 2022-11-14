FROM alpine:3.16.3

WORKDIR /Documents/sunduck/url-shorter/

COPY .env .
COPY urler .

EXPOSE 3333
CMD urler