# FROM golang:alpine as build
# WORKDIR /var/local
# COPY . .
# RUN go build ./cmd/jd/

FROM alpine:latest
WORKDIR /var/local

# COPY --from=build /var/local/jd .
COPY ./jd .
COPY ./dependencies ./dependencies

ENV LOG_LEVEL="INFO"
ENV MODULE_NAME=jdcrawler
ENV CHROME_BINARY_PATH="/var/local/dependencies/chrome-linux/chrome"
ENV CHROMEDRIVER_PATH="/var/local/dependencies/chromedriver"
ENV SELENIUM_SERVICE_PORT=36001
ENV ITEM_CRAWL_INTERVAL=1000

CMD [ "./jd" ]