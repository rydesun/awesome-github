package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"

	"github.com/rydesun/awesome-github/exch/github"
	"github.com/rydesun/awesome-github/lib/cohttp"
)

func main() {
}

type Router struct {
	app    *fiber.App
	listen string

	html       string
	scriptPath string
	dataPath   string
}

func NewRouter(listen string) (*Router, error) {
	return &Router{
		app:    fiber.New(),
		listen: listen,
	}, nil
}

func (r *Router) Init(repoID github.RepoID, scriptPath, dataPath string) error {
	html, err := FetchHTMLReadme(repoID)
	if err != nil {
		return err
	}
	r.html = html
	r.dataPath = dataPath
	r.scriptPath = scriptPath
	return nil
}

func (r *Router) Route() {
	app := r.app
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) {
		c.Set("content-type", "text/html; charset=utf-8")
		c.Send(wrapReadme(r.html, "/data", "/js"))
	})
	app.Static("/js", r.scriptPath)
	app.Static("/data", r.dataPath)
	app.Listen(r.listen)
}

func FetchHTMLReadme(id github.RepoID) (string, error) {
	gc, err := github.NewClient(cohttp.NewClient(
		http.Client{}, 1, 3, time.Second,
		0, nil), nil, github.NewDefaultClientOption())
	if err != nil {
		return "", err
	}
	return gc.GetHTMLReadme(id)
}

func wrapReadme(readme, data_url, script_url string) string {
	return fmt.Sprintf(`
	%s
	<script src="https://cdn.jsdelivr.net/npm/luxon@1.24.1/build/global/luxon.min.js"></script>
	<script>window.data_url = "%s"</script>
	<script src="%s"></script>
	`, readme, data_url, script_url)
}
