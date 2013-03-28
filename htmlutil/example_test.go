package htmlutil_test

import "log"
import "os"

import htmlutil "."

func ExampleRender() {
	n, err := htmlutil.ParseFile("testdata/0001.html")
	if err != nil {
		log.Fatalln(err)
	}
	htmlutil.Render(os.Stdout, n)
	// Output:
	// <html>
	//    <head>
	//       <meta charset="utf-8">
	//       <title>
	//          0001
	//       </title>
	//    </head>
	//    <body>
	//       <img src="test.png" alt="test">
	//       <p>
	//          Test page.
	//       </p>
	//    </body>
	// </html>
}
