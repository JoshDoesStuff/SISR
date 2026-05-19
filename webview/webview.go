package webview

type WebView interface {
	Navigate(url string)
	SetHTML(html string)
	Eval(js string)
	Bind(name string, fn any) error
	SetVisible(visible bool)
	Visible() bool
	Resize(w, h int)
	Tick()
	Destroy()
}
