// Package htmlutil implements some html utility functions.
package htmlutil

import "bytes"
import "exp/html"
import "fmt"
import "io"
import "io/ioutil"
import "strings"

// ParseFile returns the parse tree for the HTML from the provided file.
func ParseFile(htmlPath string) (node *html.Node, err error) {
	buf, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		return nil, err
	}
	node, err = html.Parse(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return node, err
}

// Render renders a simplified version of the provided html node. It will strip
// doctype nodes, trim spaces on text nodes and insert indentation.
//
// Note: Render doesn't guarantee that the semantics of the page are preserved.
func Render(w io.Writer, n *html.Node) {
	render(w, n, 0)
}

// RenderToString renders a simplified version of the provided html node and
// returns it as a string. It will strip doctype nodes, trim spaces on text
// nodes and insert indentation.
//
// Note: RenderToString doesn't guarantee that the semantics of the page are
// preserved.
func RenderToString(n *html.Node) string {
	w := new(bytes.Buffer)
	render(w, n, 0)
	return string(w.Bytes())
}

func render(w io.Writer, n *html.Node, indent int) {
	switch n.Type {
	case html.ElementNode:
		fmt.Fprintf(w, strings.Repeat("   ", indent))
		fmt.Fprintf(w, "<%s", n.Data)
		for _, attr := range n.Attr {
			fmt.Fprintf(w, " %s='%s'", attr.Key, attr.Val)
		}
		fmt.Fprintf(w, ">\n")
	case html.TextNode:
		renderText(w, n.Data, indent)
	case html.CommentNode:
		if len(strings.TrimSpace(n.Data)) == 0 {
			break
		}
		fmt.Fprintf(w, "<!-- ")
		renderText(w, n.Data, indent)
		fmt.Fprintf(w, " -->")
	}
	for _, c := range n.Child {
		i := indent + 1
		if n.Type == html.DocumentNode {
			i = indent
		}
		render(w, c, i)
	}
	switch n.Type {
	case html.ElementNode:
		fmt.Fprintf(w, strings.Repeat("   ", indent))
		fmt.Fprintf(w, "</%s>\n", n.Data)
	}
}

func renderText(w io.Writer, data string, indent int) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		s := strings.TrimSpace(line)
		if len(s) == 0 {
			continue
		}
		fmt.Fprintf(w, strings.Repeat("   ", indent))
		fmt.Fprintf(w, "%s\n", s)
	}
}
