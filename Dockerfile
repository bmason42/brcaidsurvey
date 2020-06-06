FROM openjdk:8-jre-alpine as openapigenerate
#FROM openjdk:8 as openapigenerate
LABEL stage=openapigenerate
WORKDIR /work
COPY ./ ./
#RUN make justgenerate
RUN	mkdir tmp
RUN	java -jar third_party/tools/openapi-generator-cli.jar generate -g go-gin-server --package-name v1 -i api/openapi-1.yaml -Dmodels -o tmp

FROM golang:1.12.0 as build

WORKDIR /build

RUN mkdir client & mkdir app
COPY ./ ./


RUN rm -r -f pkg/generated & mkdir -p pkg/generated/v1
COPY --from=openapigenerate /work/tmp/go/* pkg/generated/v1/

RUN make buildlinux

FROM scratch
#RUN mkdir -p third_party/swaggerui & mkdir -p api
COPY --from=build /build/out/opencrisisline2_x64linux ./opencrisisline2_x64linux
COPY --from=build /build/third_party/swaggerui/* ./third_party/swaggerui/
COPY --from=build /build/web/ ./web/
COPY --from=build /build/api/* ./api/
ENTRYPOINT ["./opencrisisline2_x64linux"]
