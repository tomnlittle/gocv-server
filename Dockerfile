FROM denismakogon/gocv-alpine:4.0.1-buildstage as build-stage

RUN go get -u -d gocv.io/x/gocv
RUN cd $GOPATH/src/gocv.io/x/gocv && go build -o $GOPATH/bin/gocv-version ./cmd/version/main.go

# ADD Gopkg.* /go/src/github.com/tomnlittle/gocv-server
# RUN dep ensure --vendor-only
ADD . /go/src/github.com/tomnlittle/gocv-server
RUN go install /go/src/github.com/tomnlittle/gocv-server

FROM denismakogon/gocv-alpine:4.0.1-runtime

COPY --from=build-stage /go/bin/gocv-version /gocv-version
COPY --from=build-stage /go/bin/gocv-server /go/bin/gocv-server
