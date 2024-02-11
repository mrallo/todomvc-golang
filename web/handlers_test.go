package web

import (
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/todomvc-golang/todo"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var templ = template.Must(template.New("index").Parse("<p>{{.Items}}</p>"))

var (
	idZero = todo.MustNewItemId("0")
	idOne  = todo.MustNewItemId("1")
	idTwo  = todo.MustNewItemId("2")
)

type ListFinderStub struct {
	model *todo.List
}

func (l ListFinderStub) FindList() (*todo.List, error) {
	return l.model, nil
}

func Test_indexHandler_ok(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	testFinder := ListFinderStub{todo.NewList()}

	IndexHandler(templ, testFinder).ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "<p>[]</p>", w.Body.String())
}

func Test_indexHandler_unexpectedPath(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/foo", nil)
	testFinder := ListFinderStub{todo.NewList()}

	IndexHandler(templ, testFinder).ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "Not found\n", w.Body.String())
}

func Test_indexHandler_editItem(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/?edit=3", nil)
	templ := template.Must(template.New("index").Parse("<p>{{.EditingItemId}}</p>"))
	testFinder := ListFinderStub{todo.NewList()}

	IndexHandler(templ, testFinder).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<p>3</p>", w.Body.String())
}

func Test_indexHandler_editItemNotPassed(t *testing.T) {
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil)
	templ := template.Must(template.New("index").Parse("<p>{{.EditingItemId}}</p>"))
	testFinder := ListFinderStub{todo.NewList()}

	IndexHandler(templ, testFinder).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<p></p>", w.Body.String())
}

func Test_editHandler_ok(t *testing.T) {
	assert := assert.New(t)
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader("todoItemId=0&todoItemTitle=bar"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	model := todo.NewList()
	model.Add("foo")
	templ := template.Must(template.New("index").Parse("<p>{{len .Items}}</p>"))

	EditHandler(templ, model).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("<p>1</p>", w.Body.String())
	assert.Equal("bar", model.Items[idZero].Title)
}

func Test_editHandler_textIsEmpty(t *testing.T) {
	assert := assert.New(t)
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader("todoItemId=0&todoItemTitle="))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	model := todo.NewList()
	model.Add("foo")

	EditHandler(templ, model).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("<p>[]</p>", w.Body.String())
	assert.Equal(0, len(model.Items))
}

type DestroyerMock struct {
	ids []todo.ItemId
}

func (d *DestroyerMock) Destroy(id todo.ItemId) error {
	d.ids = append(d.ids, id)
	return nil
}

func Test_destroyHandler_ok(t *testing.T) {
	assert := assert.New(t)
	w, r := httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader("todoItemId=123"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testFinder := ListFinderStub{todo.NewList()}
	destroyer := &DestroyerMock{}

	DestroyHandler(templ, testFinder, destroyer).ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("<p>[]</p>", w.Body.String())
	assert.Contains(destroyer.ids, todo.MustNewItemId("123"))
}
