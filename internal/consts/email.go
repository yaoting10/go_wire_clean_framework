package consts

import (
	"github.com/gophero/goal/stringx"
	"strings"
)

type Email struct {
}

func BlurAppEmail(email string, start int, end int, num int) string {
	p := strings.IndexByte(email, '@')
	if p < 0 {
		return email
	}
	name := email[0:p]
	domain := email[p:]

	if len(name) < start {
		start = len(name)
	}
	if len(name) < 5 {
		b := strings.Builder{}
		prev := name[0:start]
		suf := name[start:len(name)]
		b.WriteString(prev)
		b.WriteString("****")
		b.WriteString(suf)
		name = b.String()
	}
	blur := stringx.Blur(name, start, len(name)-end, "*", num)
	b := strings.Builder{}
	b.WriteString(blur)
	b.WriteString(domain)
	return b.String()
}
