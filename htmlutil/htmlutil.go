// Package htmlutil implements some html utility functions.
package htmlutil

import "bytes"
import "fmt"
import "io"
import "os"
import "strings"

import "code.google.com/p/go.net/html"

// ParseFile parses the provided HTML file and returns an HTML node.
func ParseFile(htmlPath string) (n *html.Node, err error) {
	f, err := os.Open(htmlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	n, err = html.Parse(f)
	if err != nil {
		return nil, err
	}
	return n, err
}

// Render renders a simplified version of the provided HTML node. It will strip
// doctype nodes, trim spaces on text nodes and insert indentation.
//
// Note: Render doesn't guarantee that the semantics of the page are preserved.
func Render(w io.Writer, n *html.Node) {
	render(w, n, 0)
}

// RenderToString renders a simplified version of the provided HTML node and
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
			fmt.Fprintf(w, ` %s="%s"`, attr.Key, html.EscapeString(attr.Val))
		}
		fmt.Fprintf(w, ">\n")
	case html.TextNode:
		renderText(w, n.Data, indent)
	case html.CommentNode:
		if len(strings.TrimSpace(n.Data)) == 0 {
			// skip empty comments.
			break
		}
		fmt.Fprintf(w, "<!--\n")
		renderText(w, n.Data, indent)
		fmt.Fprintf(w, "-->\n")
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		i := indent + 1
		if n.Type == html.DocumentNode {
			// don't increase indentation after void elements.
			i = indent
		}
		render(w, c, i)
	}
	if n.Type == html.ElementNode && !voidElements[n.Data] {
		// close non-void elements
		fmt.Fprintf(w, strings.Repeat("   ", indent))
		fmt.Fprintf(w, "</%s>\n", n.Data)
	}
}

func renderText(w io.Writer, data string, indent int) {
	lines := strings.Split(html.EscapeString(data), "\n")
	for _, line := range lines {
		s := strings.TrimSpace(line)
		if len(s) == 0 {
			continue
		}
		fmt.Fprintf(w, strings.Repeat("   ", indent))
		fmt.Fprintf(w, "%s\n", s)
	}
}

// Section 12.1.2, "Elements", gives this list of void elements. Void elements
// are those that can't have any contents.
var voidElements = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}
