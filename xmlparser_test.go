package xmlparser

import (
	"bufio"
	"os"
	"testing"
)

func getparser(prop string) *XMLParser {

	return getparserFile("sample.xml", prop)
}

func getparserFile(filename, prop string) *XMLParser {

	file, _ := os.Open(filename)

	br := bufio.NewReader(file)

	p := NewXMLParser(br, prop)

	return p

}

func TestBasics(t *testing.T) {

	p := getparser("tag1")

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}
	if len(results) != 2 {
		panic("Test failed result must be 2")
	}

	if len(results[0].Childs) != 4 || len(results[1].Childs) != 4 {
		panic("Test failed")
	}
	// result 1
	if results[0].Attrs["att1"] != "<att0>" || results[0].Attrs["att2"] != "att0" {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][0].Attrs["att1"] != "att0" {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][0].InnerText != "InnerText110" {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][1].InnerText != "InnerText111" {
		panic("Test failed")
	}

	if results[0].Childs["tag12"][0].Attrs["att1"] != "att0" {
		panic("Test failed")
	}

	if results[0].Childs["tag12"][0].InnerText != "" {
		panic("Test failed")
	}

	if results[0].Childs["tag13"][0].Attrs != nil && results[0].Childs["tag13"][0].InnerText != "InnerText13" {
		panic("Test failed")
	}

	if results[0].Childs["tag14"][0].Attrs != nil && results[0].Childs["tag14"][0].InnerText != "" {
		panic("Test failed")
	}

	//result 2
	if results[1].Attrs["att1"] != "<att1>" || results[1].Attrs["att2"] != "att1" {
		panic("Test failed")
	}

	if results[1].Childs["tag11"][0].Attrs["att1"] != "att1" {
		panic("Test failed")
	}

	if results[1].Childs["tag11"][0].InnerText != "InnerText2" {
		panic("Test failed")
	}

	if results[1].Childs["tag12"][0].Attrs["att1"] != "att1" {
		panic("Test failed")
	}

	if results[1].Childs["tag12"][0].InnerText != "" {
		panic("Test failed")
	}
	if results[1].Childs["tag13"][0].Attrs != nil && results[1].Childs["tag13"][0].InnerText != "InnerText213" {
		panic("Test failed")
	}

	if results[1].Childs["tag14"][0].Attrs != nil && results[1].Childs["tag14"][0].InnerText != "" {
		panic("Test failed")
	}

}

func TestTagWithNoChild(t *testing.T) {

	p := getparser("tag2")

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}
	if len(results) != 2 {
		panic("Test failed")
	}
	if results[0].Childs != nil || results[1].Childs != nil {
		panic("Test failed")
	}
	if results[0].Attrs["att1"] != "testattr<" || results[1].Attrs["att1"] != "testattr<2" {
		panic("Test failed")
	}
	// with inner text
	p = getparser("tag3")

	results = results[:0]
	for xml := range p.Stream() {
		results = append(results, xml)
	}

	if len(results) != 2 {
		panic("Test failed")
	}
	if results[0].Childs != nil || results[1].Childs != nil {
		panic("Test failed")
	}

	if results[0].Attrs != nil || results[0].InnerText != "tag31" {
		panic("Test failed")
	}

	if results[1].Attrs["att1"] != "testattr<2" || results[1].InnerText != "tag32 " {
		panic("Test failed")
	}

}

func TestTagWithSpaceAndSkipOutElement(t *testing.T) {

	p := getparser("tag4").SkipElements([]string{"skipOutsideTag"}).SkipOuterElements()

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}

	if len(results) != 1 {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][0].Attrs["att1"] != "att0 " {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][0].InnerText != "InnerText0 " {
		panic("Test failed")
	}

}

func TestQuote(t *testing.T) {

	p := getparser("quotetest")

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}

	if len(results) != 1 {
		panic("Test failed")
	}

	if results[0].Attrs["att1"] != "test" || results[0].Attrs["att2"] != "test\"" || results[0].Attrs["att3"] != "test'" {
		panic("Test failed")
	}

}

func TestSkip(t *testing.T) {

	p := getparser("tag1").SkipElements([]string{"tag11", "tag13"})

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}

	if len(results[0].Childs) != 2 {
		panic("Test failed")
	}

	if len(results[1].Childs) != 2 {
		panic("Test failed")
	}

	if results[0].Childs["tag11"] != nil {
		panic("Test failed")
	}

	if results[0].Childs["tag13"] != nil {
		panic("Test failed")
	}

	if results[1].Childs["tag11"] != nil {
		panic("Test failed")
	}

	if results[1].Childs["tag13"] != nil {
		panic("Test failed")
	}

}

func TestError(t *testing.T) {

	p := getparserFile("error.xml", "tag1")

	for xml := range p.Stream() {
		if xml.Err == nil {
			panic("It must give error")
		}
	}

}
func TestGetAllNodes(t *testing.T) {
	p := getparser("examples")
	for xml := range p.Stream() {
		nodes := xml.GetAllNodes("father.son.grandson")
		if len(nodes) != 8 {
			t.Errorf("Lenght of xml.GetAllNodes is not the expected \n\t Expected: %d \n\t Found: %d", 8, len(nodes))
		} else {
			values := []string{"grandson111", "grandson112", "grandson121", "grandson122", "grandson131", "grandson132", "grandson211", "grandson212"}
			for i, node := range nodes {
				if node.GetValue(".") != values[i] {
					t.Errorf("The value of the grandson %d doesn´t match with the expected \n\t Expected: %s \n\t Found: %s", i, values[i], node.GetValue("."))
				}
			}
		}
	}
}

type testStringHelper struct {
	path     string
	expected string
}
type testIntHelper struct {
	path     string
	expected int
}

func TestGetValue(t *testing.T) {
	var found string
	p := getparser("examples")
	testHelper := []testStringHelper{
		{
			path:     "@inittag",
			expected: "initial_attr",
		},
		{
			path:     "tag1.tag11",
			expected: "InnerText110",
		},
		{
			path:     "tag1.tag11[1]",
			expected: "InnerText111",
		},
		{
			path:     "tag1[1].tag11",
			expected: "InnerText2",
		},
		{
			path:     "tag1[10].tag11",
			expected: "",
		},
		{
			path:     "tag1.tag11[10]",
			expected: "",
		},
		{
			path:     "tag1.tag12@att1",
			expected: "att0",
		},
		{
			path:     "tag1[1].tag12@att1",
			expected: "att1",
		},
		{
			path:     "tag1[1].tag12@missingatt",
			expected: "",
		},
		{
			path:     "missingtag.tag12.tag13",
			expected: "",
		},
		{
			path:     "tag1[1].tag12.missingtag@att1",
			expected: "",
		},
	}
	for xml := range p.Stream() {
		for _, testH := range testHelper {
			found = xml.GetValue(testH.path)
			if found != testH.expected {
				t.Errorf("%s doesn´t match with expected \n\t Expected: %s \n\t Found: %s", testH.path, testH.expected, found)
			}
		}
		node := xml.GetNode("tag1[1].tag13")
		found = node.GetValue(".")
		if found != "InnerText213" {
			t.Errorf("tag1[1]>tag13 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "InnerText213", found)
		}
	}

}
func TestGetValueNumeric(t *testing.T) {
	var f float64
	testHelper := []testIntHelper{
		{
			path:     "numeric.int",
			expected: 8,
		},
		{
			path:     "numeric.int[1]",
			expected: 18,
		},
		{
			path:     "numeric.int[2]",
			expected: 0,
		},
		{
			path:     "numeric.int[2]@realInt",
			expected: 9,
		},
	}
	p := getparser("examples")
	for xml := range p.Stream() {
		for _, testH := range testHelper {
			found := xml.GetValueInt(testH.path)
			if found != testH.expected {
				t.Errorf("%s doesn´t match with expected \n\t Expected: %d \n\t Found: %d", testH.path, testH.expected, found)
			}
		}
		f = xml.GetValueF64("numeric.float")
		if f != 39.9 {
			t.Errorf("numeric.float doesn´t match with expected \n\t Expected: %f \n\t Found: %f", 39.9, f)
		}
	}
}
func TestGetValueDeep(t *testing.T) {
	p := getparser("examples")
	for xml := range p.Stream() {
		i := xml.GetValueIntDeep("numericDeep.deep.int")
		if i != 2 {
			t.Errorf("numericDeep.deep.int doesn´t match with expected \n\t Expected: %d \n\t Found: %d", 2, i)
		}
		f := xml.GetValueF64Deep("numericDeep.deep.float")
		if f != 1.2 {
			t.Errorf("numericDeep.deep.float doesn´t match with expected \n\t Expected: %f \n\t Found: %f", 1.2, f)
		}
		v := xml.GetValueDeep("father.son.grandson")
		if v != "grandson111" {
			t.Errorf("father.son.grandson Deep search doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "grandson111", v)
		}
		v = xml.GetValueDeep("numericDeep.deep.float")
		if v != "1.2" {
			t.Errorf("numericDeep.deep.float Deep search doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "1.2", v)
		}
	}
}

func TestCData(t *testing.T) {
	testHelper := []string{
		"OB Fees of incl for CARD FEE FDA may be applied for traveler 1.",
		"OB Fees of incl for CARD FEE FCA may be applied for traveler 1.",
		"Qantas Frequent Flyer 1 could earn 1200 Qantas Points and 20 Status Credit for this booking. <a href=\"https://www.qantas.com/fflyer/dyn/program/terms\" target=\"_blank\">Terms and conditions apply.</a>",
	}
	p := getparserFile("cdata_test.xml", "Root")
	for xml := range p.Stream() {
		warnings := xml.GetNodes("Warnings.Warning")
		for i, w := range warnings {
			expected := testHelper[i]
			found := w.GetValue(".")
			if expected != found {
				t.Errorf("Warning[%d] Search doesn´t match with expected \n\t Expected: %s \n\t Found: %s", i, expected, found)
			}
		}
		v := xml.GetValue("Response.ReShoppingResponseID.ResponseID")
		if v != "XXXXXXXX" {
			t.Errorf("Response.ReShoppingResponseID.ResponseID Search doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "XXXXXXXX", v)
		}
	}
}

func Benchmark1(b *testing.B) {

	for n := 0; n < b.N; n++ {
		p := getparser("tag4").SkipElements([]string{"skipOutsideTag"}).SkipOuterElements()
		for xml := range p.Stream() {
			nothing(xml)
		}
	}
}

func Benchmark2(b *testing.B) {

	for n := 0; n < b.N; n++ {
		p := getparser("tag4")
		for xml := range p.Stream() {
			nothing(xml)
		}
	}

}

func nothing(...interface{}) {
}
