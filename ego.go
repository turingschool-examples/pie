package pie

import (
	"fmt"
	"io"
)

//line index.ego:1
func TableIndex(w io.Writer, tables []*Table) error {
//line index.ego:2
	_, _ = fmt.Fprintf(w, "\n\n<html>\n<head>\n  <title>pie</title>\n</head>\n\n<body>\n\t<h1>Tables</h1>\n\n\t<ul>\n\t\t")
//line index.ego:12
	for _, t := range tables {
//line index.ego:13
		_, _ = fmt.Fprintf(w, "\n\t\t\t<li>")
//line index.ego:13
		_, _ = fmt.Fprintf(w, "%v", t.Name)
//line index.ego:13
		_, _ = fmt.Fprintf(w, "</li>\n\t\t")
//line index.ego:14
	}
//line index.ego:15
	_, _ = fmt.Fprintf(w, "\n\t</ul>\n</body>\n</html>\n\n")
	return nil
}
