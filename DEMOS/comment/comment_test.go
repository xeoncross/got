package comment

import (
	"testing"
)

/* Golang requires loading child templates from the parent. This is backwards
when trying to write a page and specify (from that page) which layout to use.
Using a template comment, we can have a page define the layout it inherits
from. */

var CorrectOutput = "layout: John: home: sidebar: John Footer"

func TestOrder(t *testing.T) {

	_, err := Load("templates", ".html")

	if err != nil {
		t.Error(err)
	}
}
