package htmlutil_test

import "log"
import "os"

import "github.com/mewkiz/pkg/htmlutil"

func ExampleRender() {
	n, err := htmlutil.ParseFile("testdata/0001.html")
	if err != nil {
		log.Fatalln(err)
	}
	htmlutil.Render(os.Stdout, n)
	// Output:
	// <html>
	//    <head>
	//       <meta charset='utf-8'>
	//       </meta>
	//       <title>
	//          0001
	//       </title>
	//    </head>
	//    <body>
	//       <p>
	//          Test page.
	//       </p>
	//    </body>
	// </html>
}
