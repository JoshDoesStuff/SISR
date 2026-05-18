#include "webview_linux.h"
#include <stdlib.h>

static int g_gtk_initialized = 0;

void gtk_init_once(void) {
    if (!g_gtk_initialized) {
        gtk_init();
        g_gtk_initialized = 1;
    }
}

WebViewLinux *webview_create(unsigned long sdl_xwindow, int width, int height) {
    WebViewLinux *wv = calloc(1, sizeof(WebViewLinux));
    if (!wv) {
        return NULL;
    }

    wv->window = gtk_window_new();
    gtk_window_set_decorated(GTK_WINDOW(wv->window), FALSE);
    gtk_window_set_default_size(GTK_WINDOW(wv->window), width, height);


    GtkCssProvider *css = gtk_css_provider_new();
    gtk_css_provider_load_from_string(css,
        "window.csd { margin: 0; padding: 0; box-shadow: none; border-radius: 0; }"
        "window     { background-color: transparent; }");
    gtk_style_context_add_provider_for_display(
        gdk_display_get_default(),
        GTK_STYLE_PROVIDER(css),
        GTK_STYLE_PROVIDER_PRIORITY_APPLICATION);
    g_object_unref(css);

    wv->webview = WEBKIT_WEB_VIEW(webkit_web_view_new());

    GdkRGBA transparent = {0.0, 0.0, 0.0, 0.0};
    webkit_web_view_set_background_color(wv->webview, &transparent);

    gtk_window_set_child(GTK_WINDOW(wv->window), GTK_WIDGET(wv->webview));

    gtk_widget_realize(wv->window);

    GdkDisplay *gdk_disp = gdk_display_get_default();
    Display    *x11_disp = gdk_x11_display_get_xdisplay(GDK_X11_DISPLAY(gdk_disp));
    GdkSurface *surface  = gtk_native_get_surface(GTK_NATIVE(wv->window));

    wv->gtk_xid = gdk_x11_surface_get_xid(GDK_X11_SURFACE(surface));

    XSetWindowAttributes xattrs;
    xattrs.override_redirect = True;
    XChangeWindowAttributes(x11_disp, wv->gtk_xid, CWOverrideRedirect, &xattrs);
    gdk_display_flush(gdk_disp);

    XReparentWindow(x11_disp, wv->gtk_xid, (Window)sdl_xwindow, 0, 0);
    XResizeWindow(x11_disp, wv->gtk_xid, (unsigned)width, (unsigned)height);
    XFlush(x11_disp);

    gtk_window_present(GTK_WINDOW(wv->window));
    gdk_display_flush(gdk_disp);

    return wv;
}

void webview_navigate(WebViewLinux *wv, const char *url) {
    webkit_web_view_load_uri(wv->webview, url);
}

void webview_set_html(WebViewLinux *wv, const char *html) {
    webkit_web_view_load_html(wv->webview, html, NULL);
}

void webview_eval(WebViewLinux *wv, const char *js) {
    webkit_web_view_evaluate_javascript(wv->webview, js, -1, NULL, NULL, NULL, NULL, NULL);
}

void webview_resize(WebViewLinux *wv, int width, int height) {
    gtk_window_set_default_size(GTK_WINDOW(wv->window), width, height);
    Display *x11_disp = gdk_x11_display_get_xdisplay(GDK_X11_DISPLAY(gdk_display_get_default()));
    XResizeWindow(x11_disp, wv->gtk_xid, (unsigned)width, (unsigned)height);
    XFlush(x11_disp);
}

void webview_tick(void) {
    while (g_main_context_iteration(NULL, FALSE));
}

void webview_destroy(WebViewLinux *wv) {
    if (!wv) {
        return;
    }
    if (wv->window) {
        gtk_window_destroy(GTK_WINDOW(wv->window));
    }
    free(wv);
}
