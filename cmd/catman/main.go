package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/AubSs/fasthttplogger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
)

type Entity struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   map[string]string
	Spec       map[string]string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CatalogInfoHandler(ctx *fasthttp.RequestCtx) {
	metadata := map[string]string{
		"name":        fmt.Sprintf("%s", ctx.UserValue("repo")),
		"description": time.Now().Format(time.RFC3339)}

	spec := map[string]string{
		"type":      "website",
		"lifecycle": "experimental",
		"owner":     RandStringRunes(10),
	}
	d, err := yaml.Marshal(Entity{
		ApiVersion: "backstage.io/v1alpha1",
		Kind:       "Component",
		Metadata:   metadata,
		Spec:       spec,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	ctx.Write(d)
}

func LocationsHandler(ctx *fasthttp.RequestCtx, baseUrl string) {
	locationCount, err := strconv.Atoi(fmt.Sprintf("%s", ctx.UserValue("count")))
	if err != nil {
		ctx.WriteString("failed to parse int")
		return
	}
	var result string
	for i := 1; i <= locationCount; i++ {
		location := Entity{
			ApiVersion: "backstage.io/v1alpha1",
			Kind:       "Location",
			Metadata: map[string]string{
				"name":        fmt.Sprintf("Location-%d", i),
				"description": time.Now().Format(time.RFC3339)},
			Spec: map[string]string{
				"type":   "url",
				"target": fmt.Sprintf("%s/foo/bar-%d/blob/master/catalog-info.yaml", baseUrl, i),
			},
		}
		d, err := yaml.Marshal(location)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		result = result + "---\n" + string(d)
	}

	ctx.WriteString(result)
}

func main() {
	port := flag.Int("port", 9191, "Listening port")
	baseUrl := flag.String("baseurl", "http://localhost:9191", "Base URL required for to properly create locations")
	flag.Parse()

	r := router.New()
	r.Handle("GET", "/locations/{count}/catalog-info.yaml", func(ctx *fasthttp.RequestCtx) {
		LocationsHandler(ctx, *baseUrl)
	})

	r.GET("/{org}/{repo}/blob/master/catalog-info.yaml", CatalogInfoHandler)
	log.Printf("starting catman on port %d. Base URL configured to %s", *port, *baseUrl)
	log.Fatal(fasthttp.ListenAndServe(":9191", fasthttplogger.Combined(r.Handler)))
}
