# fixie
A one gear static site generator

This project parses Markdown files in the current directory and converts them to HTML. I use it to generate the pages for my personal blog.

The idea is very similar to what [Jekyll](https://jekyllrb.com/) and [Hugo](https://gohugo.io/) do: generate HTML from Markdown files. But `fixie` is hard-coded to my needs and allows for (almost) no customizations. Hence the name: `fixie`.

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


## Layout and styling
`fixie` will process all Markdown files (`.md`) in the current directory and all its subdirectories and will generate HTML versions from the Markdown.

If there is a `layout.html` file in the current folder `fixie` will use that file as the base for all the HTML generated from the Markdown files. The `layout.html` file is expected to have a `{{CONTENT}}` tag somewhere inside of it that will be replaced with the HTML version of the Markdown for each file. You can tweak the content and structure of the `layout.html` file and update it to use your own CSS styles and whatnot.

Any Markdown files under the `./blog` are considered "blog" entries and will be used to produce an additional file `./blog/index.html` that will include a list of all the blog entries (sorted in reverse chronological order).

If the name of a file under the `./blog` directory starts with something that looks like a date (e.g. `./blog/2023-03-14/my-topic.md` or `./blog/2022-11-01-another-topic.md`) `fixie` will use that date as the creation date for the blog entry.

If the very first line of a Markdown file under the `./blog` directory is a header line (e.g. `# Hello World`) that line will be used as the title of the blog post. If no header is detected the filename will be used as the title of the blog post.

An additional `./blog/rss.xml` will be produced to list all blog entries in an RSS feed.

If is up to you to include a link to `./blog/index.html` and `./blog/rss.xml` somewhere in your `layout.html` file for them to be visible to users of the site.

