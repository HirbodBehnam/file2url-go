package shared

import (
	"file2url/database"
	"github.com/gotd/td/tg"
)

var Dispatcher = tg.NewUpdateDispatcher()
var Database database.Interface
