# fixie
A one gear static site generator

This project parses Markdown files in the current directory and converts them to HTML. I use it to generate the pages for my personal blog.

The idea is very similar to what [Jekyll](https://jekyllrb.com/) and [Hugo](https://gohugo.io/) do: generate HTML from Markdown files. But `fixie` is hard-coded to my needs and allows for (almost) no customizations. Hence the name: `fixie`.

## Getting started

You can download the `fixie` executable from [https://github.com/hectorcorrea/fixie/releases](https://github.com/hectorcorrea/fixie/releases) or via the command line with cURL:

```
$ curl -LO https://github.com/hectorcorrea/fixie/releases/latest/download/fixie
$ chmod u+x fixie
```

Once `fixie` is on your machine, make sure you have some Markdown files to process (and optionally a `layout.html`). Below are the steps to create to sample Markdown files:

```
$ echo "# Hello World" > index.md
$ echo "This is the home page" >> index.md
$ echo "and it links to the [About](about) page" >> index.md

$ echo "# About page" > about.md
$ echo "This is the about page" >> about.md

$ ls
about.md
index.md
```

Now let's run `fixie` to process these Markdown files:

```
$ ./fixie
fixie - a one gear static site generator

No layout file (./layout.html) was found
Processing .md files...
  about.md
  index.md
No blog entries (./blog/) were found
```

Notice how it reports that processed `about.md` and `index.md`. If you list your files again you should notice now the newly generated `about.html` and `index.html` files.

```
$ ls
about.md
about.html # Generated by fixie
index.md
index.html # Generated by fixie
```

At this point you can also preview your site by running `./fixie -server` and pointing your browser to `http://localhost:9001`

You might have noticed that `fixie` reported that no `layout.html` file was found. If we create one with the following content (notice that this is an HTML file with the special `{{CONTENT}}` token)

```
$ echo "<header>This is the header</header>" > layout.html
$ echo "<div>{{CONTENT}}</div>" >> layout.html
$ echo "<footer>This is the footer</footer>" >> layout.html
```

and rerun `./fixie -server` you'll see the header and footer in both pages (index.html and about.html)

Likewise, you might have noticed that when we ran `./fixie` it reported that "No blog entries (./blog/) were found". This is because in our example we don't have a `./blog` folder with Markdown files, but if you create one you should see how those files are also processed and recognized as blog entries.

## More information
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
