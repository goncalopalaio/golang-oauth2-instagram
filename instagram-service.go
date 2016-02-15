package main

/*based on http://pierrecaserta.com/go-oauth-facebook-github-twitter-google-plus/*/

import (
  "fmt"
  "golang.org/x/oauth2"
  "io/ioutil"
  "net/http"
  "os"
)

var (
  clientID = os.Getenv("ENV_INSTAGRAM_CLIENT_ID")
  clientSecret  = os.Getenv("ENV_INSTAGRAM_CLIENT_SECRET")
  redirectURL = os.Getenv("ENV_INSTAGRAM_REDIRECT_URL")
)
var lnConfig  = &oauth2.Config{
  ClientID:     clientID,
  ClientSecret: clientSecret,
  RedirectURL:  redirectURL,
  Scopes:       []string{"basic"},
  Endpoint:  oauth2.Endpoint{
    AuthURL: "https://api.instagram.com/oauth/authorize/",
    TokenURL: "https://api.instagram.com/oauth/access_token",
  },
}
var accessToken = ""

func TestGet(w http.ResponseWriter, r *http.Request){
  if len(accessToken)==0 {
        w.Write([]byte("Must have the access token first\n"))
  }

  resp,err := http.Get("https://api.instagram.com/v1/users/self/media/recent/?access_token="+accessToken)
  if err != nil {
        w.Write([]byte(err.Error()))
  }

  body, err := ioutil.ReadAll(resp.Body)
  fmt.Println(string(body))
  w.Write([]byte(body))
}

func Home(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html; charset=utf-8");
  url:= lnConfig.AuthCodeURL("")
  w.Write([]byte("<html><title>Golang Login Instagram</title> <body> <a href=" + url + "><button>Login Instagram!</button> </a> </body></html>"))
}

func InstagramLogin(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  fmt.Println(r.FormValue("code"))

  tok,err:=lnConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
  if err != nil{
    fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }

  fmt.Println(tok)
  accessToken = tok.AccessToken
  TestGet(w, r)
}

func main(){
  /*
  This code does not use Instagram -Enforce signed requests- option
  */
  mux:=http.NewServeMux()
  mux.HandleFunc("/",Home)
  mux.HandleFunc("/redirect",InstagramLogin)
  mux.HandleFunc("/get/",TestGet)
  mux.Handle("/demo/", http.StripPrefix("/demo",http.FileServer(http.Dir("demo"))))

  http.ListenAndServe(":3030", mux)
  //for the demo i used http://localhost:3030/redirect as a redirect URI/REDIRECT_URL
  //if you use a :port on your redirect url - remember to use the same one
}
