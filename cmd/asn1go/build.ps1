$COMMIT_TAG=git rev-parse --short=10 HEAD
$VERSION="v0.1.0"
$COMMIT_TIME=git log --pretty=format:%cd --date=iso HEAD -1
$GO_VERSION=go version
go build -ldflags "-X main.Version=$VERSION -X main.CommitTag=$COMMIT_TAG -X 'main.CommitTime=$COMMIT_TIME' -X 'main.GoVersion=$GO_VERSION'" .