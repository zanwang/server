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

test:
	GO_ENV=test go test ./server -v