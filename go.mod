module github.com/freman/scantp

go 1.14

replace goftp.io/server => ../../../gitea.com/freman/goftp-server

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/jlaffaye/ftp v0.0.0-20190624084859-c1312a7102bf
	github.com/minio/minio-go/v6 v6.0.46
	github.com/stretchr/testify v1.3.0
	github.com/studio-b12/gowebdav v0.0.0-20200303150724-9380631c29a1
	goftp.io/server v0.3.3
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
)
