## xml stream parser 
xml-stream-parser is xml parser for GO. It is efficient to parse large xml data with streaming fashion. 

### Usage

```xml
<?xml version="1.0" encoding="UTF-8"?>
<bookstore>
   <book isbn="XXX">
      <title>The Iliad and The Odyssey</title>
      <price>12.95</price>
      <comments>
         <userComment rating="4">Best translation I've read.</userComment>
         <userComment rating="2">I like other versions better.</userComment>
      </comments>
   </book>
   <book isbn="YYY">
      <title>Anthology of World Literature</title>
      <price>24.95</price>
      <comments>
         <userComment rating="3">Needs more modern literature.</userComment>
         <userComment rating="4">Excellent overview of world literature.</userComment>
      </comments>
   </book>
</bookstore>
```

<b>Stream</b> over books
```go


f, _ := os.Open("input.xml")
br := bufio.NewReaderSize(f,65536)
parser := xmlparser.NewXMLParser(br, "book")

for xml := range parser.Stream() {
	fmt.Println(xml.Childs["title"][0].InnerText)
	fmt.Println(xml.Childs["comments"][0].Childs["userComment"][0].Attrs["rating"])
	fmt.Println(xml.Childs["comments"][0].Childs["userComment"][0].InnerText)
}
   
```

<b>Skip</b> tags for speed
```go
parser := xmlparser.NewXMLParser(br, "book").SkipElements([]string{"price", "comments"})
```

<b>Error</b> handlings
```go
for xml := range parser.Stream() {
   if xml.Err !=nil { 
      // handle error
   }
}
```

<b>Progress</b> of parsing
```go
// total byte read to calculate the progress of parsing
parser.TotalReadSize
```

<b>Using GetValue</b> function from a XMLElement instance:
```
value = xml.GetValue("comments.userComment")
value = xml.GetValue("comments[1].userComment[1]")
```
if you would want to get the InnerText from a node:
```
value = node.GetValue(".")
```
and never do `GetValue("")`. To get an attribute value:
```
attValue = xml.GetValue("comments[1].userComment@rating")
attValue = xml.GetValue("comments.userComment[1]@rating")
attValue = xml.GetValue("@isbn")
```
<b>Using GetNodes and GetNode</b> function from a XMLElement instance:
```
singleNode = xml.GetNode("comments.userComment")
singleNode = xml.GetNode("comments[1].userComment[1]")
```
and
```
nodeArray = xml.GetNodes("comments.userComment")
nodeArray = xml.GetNodes("comments[1].userComment")
```
<b>Using GetAllNodes</b> function from a XMLElement instance:
This function allows to get all nodes following all the tree of the passed xpath and recovering all its leafs, for example
```
f, _ := os.Open("input.xml")
br := bufio.NewReaderSize(f,65536)
parser := xmlparser.NewXMLParser(br, "bookstore") 
// notice we are getting the root node, so the next for will loop only 1 time 
for xml := range parser.Stream() {
   nodes := xml.GetAllNodes("book.comment.userComment")
}
```
this invokation will return an array of 4 elements, 2 of each _comments_ node.


If you interested check also [json parser](https://github.com/tamerh/jsparser) which works similarly
