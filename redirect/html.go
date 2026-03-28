package redirect

import (
	"fmt"
	"html"
	"sort"
	"strings"
)

// FormHTML returns an HTML form that POSTs signed fields to actionURL with UTF-8.
// actionURL is typically cfg.RedirectURL().String().
func FormHTML(actionURL string, fields map[string]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<form id="cardlink-pay" accept-charset="UTF-8" action="%s" method="POST" enctype="application/x-www-form-urlencoded">`,
		html.EscapeString(actionURL))
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := fields[k]
		fmt.Fprintf(&b, `<input type="hidden" name="%s" value="%s"/>`,
			html.EscapeString(k), html.EscapeString(v))
	}
	b.WriteString(`</form>`)
	return b.String()
}
