package gee

import (
	//"log"
	"html/template"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type (
	RouteGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		parent      *RouteGroup   // support nesting
		engine      *Engine       // all groups share a Engine instance
	}

	// 整个框架的资源都是由 Engine 统一协调的
	Engine struct {
		*RouteGroup   // 继承嵌入类型的所有属性与方法
		router        *router
		groups        []*RouteGroup      // store all groups
		htmlTemplates *template.Template // for html render: 将所有模板加载入内存
		funcMap       template.FuncMap   // for html render: 所有的自定义模板渲染函数
	}
)

func New() *Engine {
	//return &Engine{router: newRouter()}
	engine := &Engine{router: newRouter()}
	engine.RouteGroup = &RouteGroup{engine: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}

	return engine
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (group *RouteGroup) Group(prefix string) *RouteGroup {
	engine := group.engine
	newGroup := &RouteGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouteGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	//log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouteGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouteGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouteGroup) Use(middleware ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middleware...)
}

// create static handler
func (group *RouteGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolution := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolution, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// check if file exist / or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
// r.Static("/assets", "/usr/daz/blog/static")
func (group *RouteGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "*/filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

func (engine *Engine) RUN(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
