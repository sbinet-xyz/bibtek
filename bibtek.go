package bibtek

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauthsvc "google.golang.org/api/oauth2/v2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", loginHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	src, err := google.DefaultTokenSource(ctx, oauthsvc.UserinfoEmailScope)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: src,
			Base:   &urlfetch.Transport{Context: ctx},
		},
	}
	client = oauth2.NewClient(ctx, src)

	svc, err := oauthsvc.New(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ui, err := svc.Userinfo.Get().Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "bibtek: Hello, world! (user=%q,%q,%q)\n", ui.Email,
		ui.FamilyName, ui.Name)

	fmt.Fprintf(w, "<div> headers</div>\n<code>\n")
	r.Header.Write(w)
	fmt.Fprintf(w, "\n</code>\n<div>done</div>\n")

	u, err := user.CurrentOAuth(ctx) //, "https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		fmt.Fprintf(w, "<code>oauth-err: %q</code>\n", err.Error())
	}
	if u != nil {
		fmt.Fprintf(w, "<code>oauth-authorized: %q</code>\n", u.Email)
	}
}

func handlerOLD(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in</a><br>`, url)
		fmt.Fprintf(w, `X-AppEngine-User-Email: %q\n`, r.Header.Get("X-AppEngine-User-Email"))
		fmt.Fprintf(w, "<div> headers</div>\n<block>\n")
		r.Header.Write(w)
		fmt.Fprintf(w, "\n</block>\n<div>done</div>\n")
		return
	}
	fmt.Fprintf(w, "sbinet.xyz/bibtek: Hello, world! (user=%q)\n", u)
	fmt.Fprintf(w, "<div> headers</div>\n<code>\n")
	r.Header.Write(w)
	fmt.Fprintf(w, "\n</code>\n<div>done</div>\n")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in!!</a>`, url)
		return
	}
	url, _ := user.LogoutURL(ctx, "/")
	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}
