package main

import "html/template"

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
		  <a href="/logout">Logout</a>
		</nav>
	  </header>

	  <nav class="main-nav">
		<a href="/apps">Applications</a>
		<a href="/apps/create">New application</a>
		<a href="/stacks">Stacks</a>
		<a href="/stacks/create">New stack</a>
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

const stacksIndexContents = `
{{define "title"}}Stacks{{end}}
{{define "contents"}}
  <h1>Stacks</h1>
{{end}}
`

const stacksCreateContents = `
{{define "title"}}New stack{{end}}
{{define "contents"}}
  <h1>New stack</h1>

  <form action="/stacks" method="post" class="form">
	<label for="name">Name</label>
	<input type="text" name="name" placeholder="e.g. Nodejs - Three Micro" />

	<label for="base_ami">Base AMI #</label>
	<input type="text" name="base_ami" placeholder="e.g. ami-f2e4119a" />

	<label for="region">AWS Region</label>
	<select name="region">
	  <option value="us-east-1">us-east-1 (N. Virginia)</option>
	  <option value="us-west-2">us-west-2 (Oregon)</option>
	  <option value="eu-west-1">eu-west-1 (Ireland)</option>
	</select>

	<label for="vm_size">VM Size (See table below)</label>
	<select name="vm_size">
	  <option value="t2.micro">t2.micro</option>
	  <option value="t2.small">t2.small</option>
	  <option value="t2.medium">t2.medium</option>
	  <option value="m3.large">m3.large</option>
	  <option value="r3.large">r3.large</option>
	  <option value="m3.xlarge">m3.xlarge</option>
	  <option value="c4.xlarge">c4.xlarge</option>
	</select>

	<label for="server_count">Server count</label>
	<select name="server_count">
	  <option value="1">1</option>
	  <option value="2">2</option>
	  <option value="3">3</option>
	  <option value="5">5</option>
	  <option value="10">10</option>
	  <option value="20">20</option>
	  <option value="30">30</option>
	  <option value="50">50</option>
	  <option value="200">200</option>
	  <option value="1000">9001</option>
	</select>

	<footer>
	  <button type="submit" class="btn btn-save">Create</button>
	</footer>
  </form>

  <h3>EC2 VM reference</h3>

  <table class="table">
	<tr>
	  <th>Name</th>
	  <th>vCPU</th>
	  <th>ECU</th>
	  <th>Memory</th>
	  <th>Disk</th>
	  <th>$/h</th>
	  <th>$/month</th>
	</tr>
	<tr>
	  <td>t2.micro</td>
	  <td>1</td>
	  <td>~</td>
	  <td>1</td>
	  <td>EBS</td>
	  <td>$0.013</td>
	  <td>$ 9.75</td>
	</tr>
	<tr>
	  <td>t2.small</td>
	  <td>1</td>
	  <td>~</td>
	  <td>2</td>
	  <td>EBS</td>
	  <td>$0.026</td>
	  <td>$ 19.50</td>
	</tr>
	<tr>
	  <td>t2.medium</td>
	  <td>2</td>
	  <td>~</td>
	  <td>4</td>
	  <td>EBS</td>
	  <td>$0.052</td>
	  <td>$ 39.00</td>
	</tr>
	<tr>
	  <td>m3.large</td>
	  <td>2</td>
	  <td>6.5</td>
	  <td>7.5</td>
	  <td>32 SSD</td>
	  <td>$0.140</td>
	  <td>$105.00</td>
	</tr>
	<tr>
	  <td>r3.large</td>
	  <td>2</td>
	  <td>6.5</td>
	  <td>15</td>
	  <td>32 SSD</td>
	  <td>$0.175</td>
	  <td>$131.25</td>
	</tr>
	<tr>
	  <td>m3.xlarge</td>
	  <td>4</td>
	  <td>13</td>
	  <td>15</td>
	  <td>2x 40 SSD</td>
	  <td>$0.280</td>
	  <td>$210.00</td>
	</tr>
	<tr>
	  <td>c4.xlarge</td>
	  <td>4</td>
	  <td>16</td>
	  <td>7.5</td>
	  <td>EBS</td>
	  <td>$0.232</td>
	  <td>$175.00</td>
	</tr>
  </table>
{{end}}
`

func fillTemplateMap() *TemplateMap {
	tm := TemplateMap{}
	tm["applications_index"] = template.Must(loadLayoutTemplate().Parse(applicationsIndexContents))
	tm["stacks_index"] = template.Must(loadLayoutTemplate().Parse(stacksIndexContents))
	tm["stacks_create"] = template.Must(loadLayoutTemplate().Parse(stacksCreateContents))
	return &tm
}
