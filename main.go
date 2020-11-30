package main

import (
	"alphaflow/models"
	"log"
	"net/http"
	"time"

	// database driver
	_ "github.com/mattn/go-sqlite3"

	// attache
	"github.com/attache/attache"
)

type AlphaFlow struct {
	// required
	attache.BaseContext

	// capabilities
	attache.DefaultEnvironment // loads environment variables from  the file pointed to by $ENV_FILE, defaults to secret/dev.env
	attache.DefaultDB          // connects to a database using $DB_DRIVER and $DB_DSN

	User *models.User
}

func (c *AlphaFlow) Init(w http.ResponseWriter, r *http.Request) {
	/* TODO: initialize context */
}

// GET /health
func (c *AlphaFlow) GET_Health() {
	attache.RenderJSON(
		c.ResponseWriter(),
		map[string]interface{}{
			"time": time.Now().Unix(),
			"ok":   true,
		},
	)
}

func main() {
	// bootstrap application for context type Alphaflow
	app, err := attache.Bootstrap(&AlphaFlow{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run(":5000"))
}
