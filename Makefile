deps:
	go get code.google.com/p/go.crypto/bcrypt
	go get github.com/asaskevich/govalidator
	go get github.com/dchest/uniuri
	go get github.com/gin-gonic/contrib/static
	go get github.com/gin-gonic/gin
	go get github.com/mholt/binding
	go get gopkg.in/unrolled/render.v1
	go get gopkg.in/yaml.v1
	go get github.com/mailgun/mailgun-go
	go get github.com/huandu/facebook
	go get github.com/golang/oauth2
	go get github.com/go-sql-driver/mysql
	go get github.com/smartystreets/goconvey
	go get bitbucket.org/liamstask/goose/cmd/goose
	go get github.com/ziutek/mymysql/godrv
	go get github.com/jinzhu/gorm

install: deps

test: export GO_ENV=test
test:
	go test -v