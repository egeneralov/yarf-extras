package view

import (
	"bytes"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/yarf-framework/yarf"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

// View is a yarf.Resource object that renders templates based on the request path.
type View struct {
	yarf.Resource

	// Public files root path
	Public string

	// Views root path
	Views string

	// View components root path
	Components string
	Debug      bool

	// Cache in-memory storage
	cache map[string]*bytes.Buffer

	sync.RWMutex
}

// New takes 3 strings for public, views and components paths respectively.
// Returns a new *View object properly configured.
func New(p, v, c string) *View {
	return &View{
		Public:     p,
		Views:      v,
		Components: c,
		cache:      make(map[string]*bytes.Buffer),
	}
}

// Get implements yarf.Resource.Get() method to write a public file or a view to the response.
func (v *View) Get(c *yarf.Context) error {
	// Try to render public file
	// Construct path
	p := path.Join(v.Public, c.Request.URL.EscapedPath())

	// Check if file exists and if it's a file
	if info, err := os.Stat(p); err == nil && !info.IsDir() {
		// Set cache headers before http.ServeFile
		c.Response.Header().Set("Cache-Control", "public, max-age=86400")

		// Render static file
		http.ServeFile(c.Response, c.Request, p)
		return nil
	}

	// Fallback to view rendering
	return v.Render(c.Request.URL.EscapedPath(), c)
}

// Render writes a template to the response based on the request path.
// All components are available to use inside any template.
// If v.Debug == true, all templates and components are parsed every time,
// otherwise it will cache the first parse for all further requests.
func (v *View) Render(route string, c *yarf.Context) error {
	// Figure out template name
	tplName := path.Base(route) + ".html"
	if route == "" || route == "/" {
		tplName = "index.html"
	}

	// Check if main template exists and if it's a file
	if info, err := os.Stat(path.Join(v.Views, path.Dir(route), tplName)); err != nil || info.IsDir() {
		// Try to render custom 404 page
		c.Response.WriteHeader(404)

		// Render 404
		return v.Render("/404", c)
	}

	// Use cache
	if !v.Debug {
		v.RLock()
		if tpl, ok := v.cache[route]; ok {
			// Unlock after read
			v.RUnlock()

			// Render cache
			_, err := c.Response.Write(tpl.Bytes())
			if err != nil {
				// Log error to default output
				log.Println("ERROR: " + err.Error())

				return yarf.ErrorNotFound()
			}

			return nil
		}

		// Unlock if not found
		v.RUnlock()
	}

	// Init filenames
	tpls := make([]string, 0)

	// Load component files
	err := v.loadComponents(v.Components, &tpls)
	if err != nil {
		return yarf.ErrorNotFound()
	}

	// Add view template at the end
	tpls = append(tpls, path.Join(v.Views, path.Dir(route), tplName))

	// Parse template and components
	tpl, err := template.ParseFiles(tpls...)
	if err != nil {
		// Log error to default output
		log.Println("ERROR: " + err.Error())

		return yarf.ErrorNotFound()
	}

	// Render template
	buff := new(bytes.Buffer)
	err = tpl.ExecuteTemplate(buff, tplName, nil)
	if err != nil {
		// Log error to default output
		log.Println("ERROR: " + err.Error())

		return yarf.ErrorNotFound()
	}

	// Minify result?
	result := new(bytes.Buffer)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)

	m.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags:      true,
		KeepWhitespace:   true,
	})

	err = m.Minify("text/html", result, buff)
	if err != nil {
		// If can't minify, just use original and log the error
		result = buff
		log.Println("ERROR: " + err.Error())
	}

	// Set cache
	if !v.Debug {
		// Save on cache
		v.Lock()
		v.cache[route] = result
		v.Unlock()

		c.Response.Header().Set("Cache-Control", "max-age=86400")
	}

	// Render response
	_, err = c.Response.Write(result.Bytes())
	if err != nil {
		// Log error to default output
		log.Println("ERROR: " + err.Error())

		return err
	}

	return nil
}

func (v *View) loadComponents(dir string, list *[]string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		// Log error to default output
		log.Println(err.Error())

		return err
	}
	for _, f := range files {
		if f.IsDir() {
			v.loadComponents(path.Join(dir, f.Name()), list)
		} else {
			*list = append(*list, path.Join(dir, f.Name()))
		}
	}

	return nil
}
