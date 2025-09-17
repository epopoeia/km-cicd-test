module github.com/leon/km-cicd-test

go 1.21

require (
	github.com/leon/km-cicd-test/common v0.0.0
	github.com/leon/km-cicd-test/module1 v0.0.0
	github.com/leon/km-cicd-test/module2 v0.0.0
)

replace (
	github.com/leon/km-cicd-test/common => ./common
	github.com/leon/km-cicd-test/module1 => ./module1
	github.com/leon/km-cicd-test/module2 => ./module2
)
