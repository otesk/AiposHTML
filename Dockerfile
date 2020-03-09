FROM golang AS build
RUN apt-get update && apt-get install git && apt-get install make
WORKDIR usr/local/html
COPY . .
RUN make dep
RUN make

FROM scratch
COPY --from=build /go/usr/local/html/bin/html /
COPY --from=build /go/usr/local/html/static/templates/ /static/templates
CMD ["/html"]
