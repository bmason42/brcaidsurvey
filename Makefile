BUILD_TAG:=$(shell date '+%Y%m%d%H%M')

OPENAPIDEF_FILE=api/openapi-1.yaml
openapicli_jar=third_party/tools/openapi-generator-cli.jar
DOCKERTAG=opencrisisline2
APPNAME=opencrisisline2
test:
	echo ${BUILD_TAG}


buildall: buildlinux buildmac

buildmac: export GOOS=darwin
buildmac: export GOARCH=amd64
buildmac: export CGO_ENABLED=0
buildmac: export GO111MODULE=on
buildmac:
	mkdir -p out
	rm -f  out/${APPNAME}_darwin
	go build -o out/${APPNAME}_darwin app/servermain.go

buildlinux:	export GOOS=linux
buildlinux: export GOARCH=amd64
buildlinux: export CGO_ENABLED=0
buildlinux: export GO111MODULE=on
buildlinux:
	mkdir -p out
	rm -f  out/${APPNAME}_x64linux
	go build -o out/${APPNAME}_x64linux app/servermain.go


build:
	docker build -t ${DOCKERTAG}:latest -f Dockerfile .

buildalldocker: build

current_dir = $(shell pwd)
unittests:

l2tests:

push:
	mkdir -p out
	echo ${BUILD_TAG}
	docker tag ${DOCKERTAG}:latest bmason42/${DOCKERTAG}:${BUILD_TAG}
	docker tag ${DOCKERTAG}:latest bmason42/${DOCKERTAG}:latest
	#docker push bmason42/${DOCKERTAG}:${BUILD_TAG}
	docker push bmason42/${DOCKERTAG}:latest
	echo Use this image name:
	echo bmason42/${DOCKERTAG}:${BUILD_TAG}
	echo "${BUILD_TAG}" >out/buildtag.txt


generate: justgenerate
	rm -r -f pkg/generated
	mkdir -p pkg/generated/v1

	cp tmp/go/* pkg/generated/v1

	rm -rf tmp

justgenerate:
	rm -r -f tmp
	mkdir tmp
	java -jar ${openapicli_jar} generate -g go-gin-server --package-name v1 -i ${OPENAPIDEF_FILE} -Dmodels -o tmp

clean:
	rm -r -f tmp
	rm -r -f pkg/generated/v1
	rm -r -f out
	rm go.sum
