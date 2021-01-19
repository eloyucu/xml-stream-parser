package xmlparser

import (
	"strconv"
	"strings"
)

type XMLElement struct {
	Name      string
	Attrs     map[string]string
	InnerText string
	Childs    map[string][]XMLElement
	Err       error
	// filled when xpath enabled
	childs    []*XMLElement
	parent    *XMLElement
	attrs     []*xmlAttr
	localName string
	prefix    string
}

type xmlAttr struct {
	name  string
	value string
}

//******************************************//
//******************************************//
//******************************************//
//***********Imperative methods*************//
//******************************************//
//******************************************//
//******************************************//

func (element *XMLElement) GetAllNodes(xpath string) []XMLElement {
	xpaths := strings.SplitN(xpath, ".", 2)
	if len(xpaths) > 1 {
		paths := xpaths[1]
		elements := []XMLElement{}
		nodes := element.GetNodes(xpaths[0])
		for _, node := range nodes {
			elements = append(elements, node.GetAllNodes(paths)...)
		}
		return elements
	}
	return element.GetNodes(xpaths[0])
}
func (element *XMLElement) GetNodes(xpath string) []XMLElement {
	var path, paths string
	xpaths := strings.SplitN(xpath, ".", 2)
	if len(xpaths) > 1 {
		paths = xpaths[1]
	}
	path = xpaths[0]
	path, index := element.pathIndex(path)
	if len(element.Childs[path]) > index {
		if paths == "" {
			return element.Childs[path]
		}
		return element.Childs[path][index].GetNodes(paths)
	}
	return []XMLElement{}
}
func (element *XMLElement) GetNode(xpath string) XMLElement {
	var index int
	nodes := element.GetNodes(xpath)
	indexes := strings.Split(xpath, ".")
	indexes = strings.Split(indexes[len(indexes)-1], "[")
	indexes = strings.Split(indexes[len(indexes)-1], "[")
	if len(indexes) == 1 {
		var err error
		index, err = strconv.Atoi(strings.Split(indexes[0], "]")[0])
		if err != nil {
			index = 0
		}
	} else {
		index = 0
	}
	if len(nodes) > index {
		return nodes[index]
	}
	return XMLElement{}
}
func (element *XMLElement) GetValueF64(xpath string) float64 {
	v := element.GetValue(xpath)
	f := 0.00
	if t, err := strconv.ParseFloat(v, 64); err == nil {
		f = t
	}
	return f
}
func (element *XMLElement) GetValueInt(xpath string) int {
	i := element.GetValueF64(xpath)
	return int(i)
}
func (element *XMLElement) GetValue(xpath string) string {
	if xpath == "." {
		return element.InnerText
	} else if xpath == "" {
		return ""
	}
	var attr string
	var node XMLElement
	xpaths := strings.SplitN(xpath, "@", 2)
	if len(xpaths) > 1 {
		attr = xpaths[1]
	}
	if xpaths[0] == "" {
		node = *element
	} else {
		node = element.GetNode(xpaths[0])
	}
	if attr == "" {
		return node.InnerText
	}
	return node.Attrs[attr]
}
func (element *XMLElement) GetValueF64Deep(xpath string) float64 {
	v := element.GetValueDeep(xpath)
	f := 0.00
	if t, err := strconv.ParseFloat(v, 64); err == nil {
		f = t
	}
	return f
}
func (element *XMLElement) GetValueIntDeep(xpath string) int {
	i := element.GetValueF64Deep(xpath)
	return int(i)
}
func (element *XMLElement) GetValueDeep(xpath string) string {
	xpaths := strings.SplitN(xpath, ".", 2)
	if len(xpaths) > 1 {
		paths := xpaths[1]
		nodes := element.GetNodes(xpaths[0])
		for _, node := range nodes {
			v := node.GetValueDeep(paths)
			if v != "" {
				return v
			}
		}
	}
	return element.GetValue(xpaths[0])
}
func (element *XMLElement) pathIndex(path string) (string, int) {
	indexes := strings.Split(path, "[")
	path = indexes[0]
	if len(indexes) > 1 {
		indexes := strings.Split(indexes[1], "]")
		index, err := strconv.Atoi(indexes[0])
		if err != nil {
			return path, 0
		}
		return path, index
	}
	return path, 0
}
