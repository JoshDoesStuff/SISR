package webview

type WebView interface {
	Navigate(url string)
	SetHTML(html string)
	Eval(js string)
	Bind(name string, fn interface{}) error
	Resize(w, h int)
	Tick()
	Destroy()
}
