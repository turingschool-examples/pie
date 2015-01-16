package pie

import (
	"fmt"
	"io"
)

//line head.ego:1
func head(w io.Writer) error {
//line head.ego:2
	_, _ = fmt.Fprintf(w, "\n\n<head>\n  <title>pie</title>\n  <script src=\"/assets/dropzone.js\"></script>\n</head>\n")
	return nil
}

//line index.ego:1
func Index(w io.Writer, tables []*Table) error {
//line index.ego:2
	_, _ = fmt.Fprintf(w, "\n\n<html>\n")
//line index.ego:4
	head(w)
//line index.ego:5
	_, _ = fmt.Fprintf(w, "\n\n<body>\n\t<h1>PIE</h1>\n\n\t<h2>Tables</h2>\n\n\t<ul>\n\t\t")
//line index.ego:12
	for _, t := range tables {
//line index.ego:13
		_, _ = fmt.Fprintf(w, "\n\t\t\t<li><a href=\"/tables/")
//line index.ego:13
		_, _ = fmt.Fprintf(w, "%v", t.Name)
//line index.ego:13
		_, _ = fmt.Fprintf(w, "\">")
//line index.ego:13
		_, _ = fmt.Fprintf(w, "%v", t.Name)
//line index.ego:13
		_, _ = fmt.Fprintf(w, "</a></li>\n\t\t")
//line index.ego:14
	}
//line index.ego:15
	_, _ = fmt.Fprintf(w, "\n\t</ul>\n\n\n\t<br/><br/>\n\t<form class=\"dropzone\" action=\"/tables\"></form>\n\n</body>\n</html>\n\n")
	return nil
}

//line show.ego:1
func TableShow(w io.Writer, t *Table, rows [][]string) error {
//line show.ego:2
	_, _ = fmt.Fprintf(w, "\n\n<html>\n<head>\n  <title>pie : ")
//line show.ego:5
	_, _ = fmt.Fprintf(w, "%v", t.Name)
//line show.ego:5
	_, _ = fmt.Fprintf(w, "</title>\n</head>\n\n<body>\n\t<h1>")
//line show.ego:9
	_, _ = fmt.Fprintf(w, "%v", t.Name)
//line show.ego:9
	_, _ = fmt.Fprintf(w, "</h1>\n\n\t<table>\n\t\t<tr>\n\t\t\t")
//line show.ego:13
	for _, c := range t.Columns {
//line show.ego:14
		_, _ = fmt.Fprintf(w, "\n\t\t\t\t<td>")
//line show.ego:14
		_, _ = fmt.Fprintf(w, "%v", c.Name)
//line show.ego:14
		_, _ = fmt.Fprintf(w, "</td>\n\t\t\t")
//line show.ego:15
	}
//line show.ego:16
	_, _ = fmt.Fprintf(w, "\n\t\t</tr>\n\n\t\t")
//line show.ego:18
	for _, row := range rows {
//line show.ego:19
		_, _ = fmt.Fprintf(w, "\n\t\t\t<tr>\n\t\t\t\t")
//line show.ego:20
		for _, value := range row {
//line show.ego:21
			_, _ = fmt.Fprintf(w, "\n\t\t\t\t\t<td>")
//line show.ego:21
			_, _ = fmt.Fprintf(w, "%v", value)
//line show.ego:21
			_, _ = fmt.Fprintf(w, "</td>\n\t\t\t\t")
//line show.ego:22
		}
//line show.ego:23
		_, _ = fmt.Fprintf(w, "\n\t\t\t</tr>\n\t\t")
//line show.ego:24
	}
//line show.ego:25
	_, _ = fmt.Fprintf(w, "\n\t</table>\n</body>\n</html>\n\n")
	return nil
}

//line visualize.ego:1
func Visualize(w io.Writer) error {
//line visualize.ego:2
	_, _ = fmt.Fprintf(w, "\n\n<html>\n")
//line visualize.ego:4
	head(w)
//line visualize.ego:5
	_, _ = fmt.Fprintf(w, "\n\n<script src=\"/assets/d3.v3.min.js\"></script>\n\n<body>\n\t<h1>Visualize</h1>\n\n\t<h2>Query</h2>\n\t<textarea cols=\"80\" rows=\"5\"></textarea>\n\n\t<hr/>\n\n\t<div id=\"chart\"></div>\n</body>\n<script src=\"/assets/visualize.js\"></script>\n\n</html>\n\n")
	return nil
}
