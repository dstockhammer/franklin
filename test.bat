REM Install Dependencies
go get github.com/streadway/amqp
go get github.com/smartystreets/goconvey
go get github.com/smartystreets/goconvey/convey

REM Start GoConvey
%GOPATH%/bin/goconvey
