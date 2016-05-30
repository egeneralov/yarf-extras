package tpl

import (
    "github.com/yarf-framework/yarf"
)

// LayoutMiddleware is a Yarf Middleware used to render fixed template files on before and after each yarf resource is executed. 
// The TplPath property should point to the root templates path where all directories and template files are centralized. 
// The TplName is used is used to build the entire files path as: [LayoutMiddleware.TplPath] + "/" + [LayoutMiddleware.TplName] + "/layout/" + [Method_Name(pre,post,end)] + ".tpl"
type LayoutMiddleware struct {
    Tpl
    yarf.Middleware
}

// NewLayoutMiddleware() constructs and sets up a new LayoutMiddleware object. 
// It uses the path and name string parameters to fill the corresponding object properties.
func NewLayoutMiddleware(path, name string) *LayoutMiddleware {
    if name == "" {
        name = "default"
    }
    
    return &LayoutMiddleware {
        Tpl: Tpl {
            TplPath: path,
            TplName: name,
            cache: make(map[string]string),
        },
    }
}

// Render reads the content of the path built using the name param together with the object properties.
// It outputs the content of the readed file or returns an error if that fails. 
func (l *LayoutMiddleware) Render(name string, c *yarf.Context) error {
    content, err := l.Cached(l.TplPath + "/" + l.TplName + "/layout/" + name + ".tpl")
    if err != nil {
        return yarf.ErrorNotFound()
    }
    
    c.Render(content)
    
    return nil
}

// PreDispatch is the default implementation and just calls l.Render("pre")
// This can be overloaded by extending LayoutMiddleware using composition and overwriting this method.
func (l *LayoutMiddleware) PreDispatch(c *yarf.Context) error {
    return l.Render("pre", c)
}

// PostDispatch is the default implementation and just calls l.Render("post")
// This can be overloaded by extending LayoutMiddleware using composition and overwriting this method.
func (l *LayoutMiddleware) PostDispatch(c *yarf.Context) error {
    return l.Render("post", c)
}

// End is the default implementation and just calls l.Render("end")
// This can be overloaded by extending LayoutMiddleware using composition and overwriting this method.
func (l *LayoutMiddleware) End(c *yarf.Context) error {
    return l.Render("end", c)
}
