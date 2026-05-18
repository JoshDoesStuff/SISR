#pragma once

#include <gtk/gtk.h>
#include <webkit/webkit.h>
#include <gdk/x11/gdkx.h>
#include <X11/Xlib.h>

typedef struct {
    GtkWidget     *window;
    WebKitWebView *webview;
    Window         gtk_xid;
} WebViewLinux;

void gtk_init_once(void);

WebViewLinux *webview_create(unsigned long sdl_xwindow, int width, int height);

void webview_navigate(WebViewLinux *wv, const char *url);
void webview_set_html(WebViewLinux *wv, const char *html);
void webview_eval(WebViewLinux *wv, const char *js);
void webview_resize(WebViewLinux *wv, int width, int height);
void webview_tick(void);
void webview_destroy(WebViewLinux *wv);
