package swagger

type Swagger struct {
	Swagger     string                `json:"swagger"`
	Info        Info                  `json:"info"`
	Paths       map[string]*PathItem  `json:"paths"`
	Definitions map[string]Definition `json:"definitions"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type PathItem struct {
	operation *Operation `json:"-"`
	method    string     `json:"-"`
	Get       *Operation `json:"get,omitempty"`
	Post      *Operation `json:"post,omitempty"`
}

type Operation struct {
	Responses  map[string]*Response `json:"responses"`
	Parameters []Parameter          `json:"parameters"`
}

type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Type        string  `json:"type"`
	Schema      *Schema `json:"schema,omitempty"`
}

type Response struct {
	Description string      `json:"description"`
	Schema      Schema      `json:"schema"`
	Examples    interface{} `json:"examples"`
}

type Schema struct {
	Ref        string           `json:"$ref"`
	Properties map[string]Items `json:"properties"`
	Required   []string         `json:"required"`
	Example    interface{}      `json:"example"`
}

type Definition struct {
	Properties map[string]Parameter `json:"properties"`
	Xml        map[string]string    `json:"xml"`
}

type Items struct {
	Type   string `json:"type"`
	Format string `json:"format"`
}

var DefaultSwagger = &Swagger{Swagger: "2.0",
	Paths:       make(map[string]*PathItem),
	Definitions: make(map[string]Definition),
}

func (s *Swagger) NewPathItem(path string, method string) *PathItem {
	pt := NewPathItem(method)
	s.Paths[path] = pt
	return pt
}

func NewPathItem(method string) *PathItem {
	op := &Operation{
		Responses:  make(map[string]*Response),
		Parameters: []Parameter{},
	}
	pt := &PathItem{method: method}
	if pt.method == "GET" {
		pt.operation = op
		pt.Get = op
	} else if pt.method == "POST" {
		pt.operation = op
		pt.Post = op
	} else {
		panic("the doc gen api only support GET and POST method")
	}
	return pt
}

func (s *Swagger) PutPathItem(path string, item *PathItem) {
	s.Paths[path] = item
}

func (pt *PathItem) AddPostParams(key string, params []Parameter) *PathItem {
	properties := make(map[string]Parameter)
	for _, param := range params {
		properties[param.Name] = param
	}
	_, ok := DefaultSwagger.Definitions[key]
	if !ok {
		DefaultSwagger.Definitions[key] = Definition{Properties: properties, Xml: map[string]string{"name": key}}
	}
	if len(pt.Post.Parameters) == 0 {
		pt.Post.Parameters = append(pt.Post.Parameters, Parameter{In: "body", Name: "body", Schema: &Schema{Ref: "#/definitions/" + key}, Required: true, Description: "暂缺"})
	}
	return pt

}

func (pt *PathItem) AddParam(param Parameter) *PathItem {
	pt.operation.Parameters = append(pt.operation.Parameters, param)
	return pt
}
