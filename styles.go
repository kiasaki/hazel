package main

import "net/http"

const stylesContents = `
html, body {
  margin: 0; padding: 0; height: 100%;
  font-family: 'Helvetica', sans-serif;
  font-size: 16px;
  background: #ecf0f1;
}

* {
  box-sizing: border-box;
}

img {
  max-width: 100%;
}

.btn {
  line-height: 20px;
  padding: 10px 15px;
  border: none;
  background: #859900;
  color: #fff;
  font-weight: bold;
  letter-spacing: 1px;
  cursor: pointer;
  display: inline-block;
}
a.btn {
  text-decoration: none;
}

div.wrapper {
  position: relative;
  width: 960px;
  margin: 0 auto;
  box-shadow: 0 0 4px rgba(0,0,0,0.25);
  background: #fff;
  height: 100%;
  min-height: 100%;
}

header.main-header {
  background: #8e44ad;
  padding: 20px;
  line-height: 30px;
  display: flex;
}
header.main-header h1 {
  flex: 1;
  margin: 0;
  color: #fff;
  font-size: 30px;
  font-weight: 300;
}
header.main-header nav {
  flex: 1;
  text-align: right;
}
header.main-header nav a {
  color: #fff;
  text-decoration: none;
}
header.main-header nav a:hover {
  text-decoration: underline;
}

nav.main-nav {
  background: #9b59b6;
  line-height: 20px;
  padding: 0 5px;
}
nav.main-nav a {
  color: #fff;
  text-decoration: none;
  display: inline-block;
  padding: 15px 10px;
  transition: background 150ms;
}
nav.main-nav a:hover {
  background: #8e44ad;
}

section.main-content {
  padding: 15px;
}
section.main-content h1 {
  margin: 0 0 15px;
}

footer.main-footer {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  background: #9b59b6;
  padding: 15px;
  line-height: 20px;
  color: #fff;
}
footer.main-footer a {
  color: #bdc3c7;
}
`

func RenderStyles(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(stylesContents))
}
