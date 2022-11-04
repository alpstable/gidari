module github.com/alpstable/gidari/e2e

go 1.19

replace github.com/alpstable/gidari => ../

replace github.com/alpstable/gmongo => /Users/prestonvasquez/Developer/gmongo // remove this

require (
	github.com/alpstable/gidari v0.0.0-20221030223101-0562042f4484
	github.com/sirupsen/logrus v1.9.0
)

require (
	github.com/alpstable/gmongo v0.0.0-20221031043101-be47f52cf05a // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.10.3 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
