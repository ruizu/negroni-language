package locale

import (
	"fmt"
	"net/http"
)

type Middleware struct {
	formName   string
	cookieName string
	locales    []string
}

var availableLocales []string = []string{
	"en-US",
	"id-ID",
}

func New() *Middleware {
	return NewCustomName("locale", "locale")
}

func NewCustomLocales(l []string) *Middleware {
	m := NewCustomName("locale", "locale")
	m.SetLocales(l)
	return m
}

func NewCustomName(f, c string) *Middleware {
	return &Middleware{formName: f, cookieName: c, locales: availableLocales}
}

func (m *Middleware) SetLocales(l []string) error {
	if len(l) < 1 {
		return fmt.Errorf("no locale supplied")
	}
	m.locales = l
	return nil
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	l := r.FormValue(m.formName)
	if l != "" {
		r.Form.Set(m.formName, m.getLocale(w, l))
		next(w, r)
		return
	}

	c, err := r.Cookie(m.cookieName)
	if err != nil {
		r.Form.Set(m.formName, m.getLocale(w, m.locales[0]))
		next(w, r)
		return
	}

	r.Form.Set(m.formName, m.getLocale(w, c.Value))
	next(w, r)
}

func (m *Middleware) getLocale(w http.ResponseWriter, l string) string {
	tl := m.locales[0]
	for _, v := range m.locales {
		if l == v {
			tl = l
			break
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:  m.cookieName,
		Value: tl,
	})
	return tl
}
