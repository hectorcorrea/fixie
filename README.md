# fixie
A one gear static site generator

This project parses Markdown files in the current directory and converts them to HTML. I use it to generate the pages for my personal blog.

The idea is very similar to what [Jekyll](https://jekyllrb.com/) and [Hugo](https://gohugo.io/) do: generate HTML from Markdown files. But `fixie` is hard-coded to my needs and allows for (almost) no customizations. Hence the name: `fixie`.


## How to use it
Check out this blog post: https://hectorcorrea.com/blog/2023-10-18/fixie


## Build from source
The code is written in Go and you'll need to have Go installed on your machine to build it.

```
git clone https://github.com/hectorcorrea/fixie
cd fixie
go build
```

Once you have the executable you can test it with the sample site provided in the repo. To process the Markdown files in the sample site run:

```
cd sample
../fixie
```

...and after that you can run `../fixie -server` to stand up a local server and browse to http://localhost:9001 to see the generated site.
