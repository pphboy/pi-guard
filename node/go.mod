module go-node

go 1.22.0

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/sys v0.18.0 // indirect

)

require pglib v0.0.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	gorm.io/driver/sqlite v1.5.5 // indirect
	gorm.io/gorm v1.25.7 // indirect
)

replace pglib => ../pglib
