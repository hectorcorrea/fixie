# fixie
A one gear blog engine

This project parses Markdown files in the current directory and converts them to HTML. I use it to generate the pages for my personal blog.

The idea is very similar to what [Jekyll](https://jekyllrb.com/) and [Hugo](https://gohugo.io/) do: generate HTML from Markdown files. But fixie is hard-coded to my needs and has almost no customization options. Hence the name fixie.

## Build from source
The code is written in Go and you'll need to have Go installed on your machine to build it.

```
git clone https://github.com/hectorcorrea/fixie
cd fixie
go build
```

Once you have the executable you can test it with the sample site. To process the Markdown files in the sample site run:

```
cd sample
../fixie
```

...and after that you can run `../fixie -server` to stand up a local server and browse to http://localhost:9001 to see the generated site.


## Layout and styling
`fixie` will process all Markdown files (`.md`) in the current directory and all its subdirectories and will generate HTML versions from the Markdown.

If there is a `layout.html` file in the current folder it will use that file as the base for all the HTML generated. The `layout.html` file is expected to have a `{{CONTENT}}` tag somewhere inside of it and that token will be replaced with the HTML version of the Markdown for each file. You can tweak the `layout.html` file as you with and use your own CSS styles and whatnot.

Markdown files under the `./blog` are considered "blog" entries and will be use to produce an additional file `./blog/index.html` that will include a list of all the blog entries sorted in reverse chronological order.

If the name of a file under the `./blog` directory starts with something that looks like a date (e.g. `./blog/2023-03-14/my-topic.md` or `./blog/2022-11-01-another-topic.md`) `fixie` will use that date as the creation date.

If the very first line of a markdown inside the `./blog` folder is a header line (e.g. `# Hello World`) that will be used as the topic of the blog post. If no header is detected the filename will be used as the topic of the blog post.

TODO: document the metadata use for the RSS feed.
