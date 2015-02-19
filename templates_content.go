package main

const layoutContents = `
{{define "layout"}}
<!DOCTYPE html>
<html>
  <head>
  <title>Hazel - {{template "title" .}}</title>
  <link rel="stylesheet" href="/styles.css" />
  </head>
  <body>
	<div class="wrapper">
	  <header class="main-header">
		<h1>Hazel</h1>
		<nav>
		  <a href="/apps">Applications</a>
		</nav>
	  </header>

	  <nav class="main-nav">
		<a href="/apps">Applications</a>
		<a href="/apps/new">New application</a>
	  </nav>

	  <section class="main-content">
		{{template "contents" .}}
	  </section>

	  <footer class="main-footer">
		<a href="http://github.com/kiasaki/hazel">Hazel</a>
		is open source software by
		<a href="http://github.com/kiasaki">kiasaki</a>
	  </footer>
	</div>
  </body>
</html>
{{end}}
`

const applicationsIndexContents = `
{{define "title"}}Applications{{end}}
{{define "contents"}}
  <h1>Applications</h1>
{{end}}
`
