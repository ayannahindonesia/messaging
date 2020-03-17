# ===== Lintas Arta =====
FROM golang:alpine  AS build-env

ARG APPNAME="messaging"
ARG ENV="staging"

#RUN adduser -D -g '' golang
#USER root

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update git gcc libc-dev;
RUN apk --no-cache add curl
#  tzdata wget gcc libc-dev make openssl py-pip;
ENV TZ=Asia/Jakarta
ENV GOOGLE_APPLICATION_CREDENTIALS=/go/src/google-services.json
RUN go get -u github.com/golang/dep/cmd/dep

RUN cd $GOPATH/src/"${APPNAME}"
RUN cp deploy/dev-config.yaml config.yaml
RUN cp deploy/dev-google-services.json google-services.json
RUN dep ensure -v
RUN go build -v -o "${APPNAME}-res"

RUN ls -alh $GOPATH/src/
RUN ls -alh $GOPATH/src/"${APPNAME}"
RUN ls -alh $GOPATH/src/"${APPNAME}"/vendor
RUN pwd

FROM alpine

WORKDIR /go/src/
ENV TZ=Asia/Jakarta
ENV GOOGLE_APPLICATION_CREDENTIALS=/go/src/google-services.json

COPY --from=build-env /go/src/messaging/messaging-res /go/src/messaging
COPY --from=build-env /go/src/messaging/deploy/dev-config.yaml /go/src/config.yaml
COPY --from=build-env /go/src/messaging/deploy/dev-google-services.json /go/src/google-services.json

RUN pwd
#ENTRYPOINT /app/asira_borrower-res
CMD ["/go/src/messaging","run"]

EXPOSE 8000

# ===== Development =====
# FROM golang:alpine

# ARG APPNAME="messaging"
# ARG ENV="dev"

# ADD . $GOPATH/src/"${APPNAME}"
# WORKDIR $GOPATH/src/"${APPNAME}"

# RUN apk add --update git gcc libc-dev tzdata;
# #  tzdata wget gcc libc-dev make openssl py-pip;
# ENV TZ=Asia/Jakarta
# ENV GOOGLE_APPLICATION_CREDENTIALS=$GOPATH/src/"${APPNAME}"/google-services.json

# RUN go get -u github.com/golang/dep/cmd/dep

# CMD if [ "${ENV}" = "dev" ] ; then \
#         cp deploy/dev-config.yaml config.yaml ; \
#         cp deploy/dev-google-services.json google-services.json ; \
#     fi \
#     && dep ensure -v \
#     && go build -v -o $GOPATH/bin/"${APPNAME}" \
#     # run app mode
#     && "${APPNAME}" run \
#     # update db structure
#     && if [ "${ENV}" = "dev"] ; then \
#         "${APPNAME}" migrate up \
#         && "${APPNAME}" seed ; \
#     fi \
#     && go test tests/*_test.go -failfast -v ;

# EXPOSE 8009