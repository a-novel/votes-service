package services_test

import (
	"fmt"
	"time"
)

var (
	fooErr = fmt.Errorf("foo")
)

var (
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 9, 0, 0, 0, time.UTC)
)
