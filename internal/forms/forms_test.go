package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}
func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")
	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error(" required field missing")
	}
}
func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	data := "wertyu"
	if !form.IsEmail(data) {
		if form.Valid() {
			t.Error("IsEmail has failed 1")
		}
	}
	data = "wertyu@dfg.com"
	if form.IsEmail(data) {
		if form.Valid() {
			t.Error("IsEmail has failed 2")
		}
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	data := "dfghj"
	if form.Has(data) {
		if form.Valid() {
			t.Error("Has has failed 1")
		}
	}

	data = "werty"
	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	if !form.Has(data) {
		if form.Valid() {
			t.Error("Has has failed 2")
		}
	}
}
func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	data := "we3"
	if !form.MinLength(data, 5) {
		if form.Valid() {
			t.Error("MinLength has failed 1")
		}
	}
	data = "welcome"
	form.MinLength(data, 5)
	if form.IsEmail(data) {
		if form.Valid() {
			t.Error("MinLength has failed 2")
		}
	}
}
