package tpl

import (
    "github.com/yarf-framework/yarf"
)

// ViewResource implements only GET Yarf resource method.
// This resourse is designed to be used in composition with any other resource that composites yarf.Resource
// It also renders the template files corresponding to the defined tpl name. 
// The TplName is used is used to build the entire files path as: [.TplPath] + "/" + [.TplName] + "/view" + [Request.URL.Path] + ".tpl"
type ViewResource struct {
    Tpl
    yarf.Resource
}

// NewLayoutMiddleware() constructs and sets up a new LayoutMiddleware object. 
// It uses the path and name string parameters to fill the corresponding object properties.
func NewViewResource(path, name string) *ViewResource {
    if name == "" {
        name = "default"
    }
    
    return &ViewResource {
        Tpl: Tpl {
            TplPath: path,
            TplName: name,
            cache: make(map[string]string),
        },
    }
}

// Render reads the content of the path built using the c.URL.Path together with the object's TplPath.
// It outputs the content of the readed file or returns an error if that fails. 
func (v *ViewResource) Render(c *yarf.Context) error {
    pre, err := v.Cached(v.TplPath + "/" + v.TplName + "/layout/pre.tpl")
    if err != nil {
        return yarf.ErrorNotFound()
    }
    
    content, err := v.Cached(v.TplPath + "/" + v.TplName + "/view" + c.Request.URL.EscapedPath() + ".tpl")
    if err != nil {
        return yarf.ErrorNotFound()
    }
    
    post, err := v.Cached(v.TplPath + "/" + v.TplName + "/layout/post.tpl")
    if err != nil {
        return yarf.ErrorNotFound()
    }
    
    end, err := v.Cached(v.TplPath + "/" + v.TplName + "/layout/end.tpl")
    if err != nil {
        return yarf.ErrorNotFound()
    }
    
    c.Render(pre)
    c.Render(content)
    c.Render(post)
    c.Render(end)
    
    return nil
}

func (v *ViewResource) Get(c *yarf.Context) error {
    return v.Render(c)
}