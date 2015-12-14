package locale

import (
	"fmt"
	"net/http"
)

type Middleware struct {
	formName   string
	cookieName string
	languages  []string
}

var availableLanguages []string = []string{
	"en-US",
	"id-ID",
}

func New() *Middleware {
	return NewCustomName("language", "language")
}

func NewCustomLanguage(l []string) *Middleware {
	m := NewCustomName("language", "language")
	m.SetLanguages(l)
	return m
}

func NewCustomName(f, c string) *Middleware {
	return &Middleware{formName: f, cookieName: c, languages: availableLanguages}
}

func (m *Middleware) SetLanguages(l []string) error {
	if len(l) < 1 {
		return fmt.Errorf("no language supplied")
	}
	m.languages = l
	return nil
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	l := r.FormValue(m.formName)
	if l != "" {
		r.Form.Set(m.formName, m.getLanguage(w, l))
		next(w, r)
		return
	}

	c, err := r.Cookie(m.cookieName)
	if err != nil {
		r.Form.Set(m.formName, m.getLanguage(w, m.languages[0]))
		next(w, r)
		return
	}

	r.Form.Set(m.formName, m.getLanguage(w, c.Value))
	next(w, r)
}

func (m *Middleware) getLanguage(w http.ResponseWriter, l string) string {
	tl := m.languages[0]
	for _, v := range m.languages {
		if l == v {
			tl = l
			break
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name: m.cookieName,
		Value: tl,
	})
	return tl
}
