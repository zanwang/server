deps:
	go get code.google.com/p/go.crypto/bcrypt
	go get github.com/asaskevich/govalidator
	go get github.com/coopernurse/gorp
	go get github.com/dchest/uniuri
	go get github.com/gin-gonic/contrib/static
	go get github.com/gin-gonic/gin
	go get github.com/mattn/go-sqlite3
	go get github.com/mholt/binding
	go get gopkg.in/unrolled/render.v1
	go get gopkg.in/yaml.v1
	go get github.com/onsi/gomega
	go get github.com/franela/goblin
	go get github.com/mailgun/mailgun-go
	go get github.com/huandu/facebook
	go get github.com/mrjones/oauth
	go get github.com/golang/oauth2

install: deps
	bower install
	npm install

test: export GO_ENV=test
test:
	go test ./tests -v