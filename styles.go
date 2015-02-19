package main

import "net/http"

const stylesContents = `
html, body {
  margin: 0; padding: 0;
  font-family: 'Courier New', sans-serif;
  font-size: 16px;
  background: #fdf6e3;
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
`

func RenderStyles(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(stylesContents))
}
