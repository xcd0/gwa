//+build js

package main

import "syscall/js"

// Document - document object
var Document = js.Global().Get("document")

// Body - document.body
var Body = Document.Get("body")

// Head - document.head
var Head = Document.Get("head")

func styles() js.Value {
	style := Document.Call("createElement", "style")

	style.Set("textContent", "html, body { height: 100%; margin: 0; display: grid; place-items: center; font-family: sans-serif }")

	return style
}

func main() {

	// Header (h1)
	h1 := Document.Call("createElement", "h1")

	h1.Set("textContent", "Congratulations! 🎉")

	// <code>

	code := Document.Call("createElement", "code")

	code.Set("textContent", "go-web-app")

	// Header (h2)

	h2 := Document.Call("createElement", "h2")

	h2.Set("textContent", "You just created a new app using ")

	h2.Call("appendChild", code)

	// Link
	a := Document.Call("createElement", "a")

	a.Set("textContent", "Star the project")

	a.Set("href", "https://github.com/talentlessguy/go-web-app")


	// App root
	root := Document.Call("createElement", "div")

	root.Call("appendChild", h1)
	root.Call("appendChild", h2)
	root.Call("appendChild", a)

	Body.Call("appendChild", root)
	Head.Call("appendChild", styles())
}
