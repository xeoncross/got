# Got
## Go Templates

The native [`html/template` engine](https://golang.org/pkg/html/template/) can handle variables, template inheritance, helper functions, and sanitizing of HTML, CSS, Javascript, and URIs. However, using `html/template` _easily_ requires a bit of boilerplate which this package provides in a minimal wrapper.

This package is for people who want to stick to closely to vanilla Go when building/using HTML templates.


## Template Functions

Template functions provide handy helpers for doing common tasks. The [Masterminds/sprig](https://github.com/Masterminds/sprig) package contains +100 helper functions (inspired by `underscore.js`) you can use to augment your templates.

If building HTML forms, or using the CSS framework Bootstrap, you might want to look at [gobuffalo/tags](https://github.com/gobuffalo/tags) for helper functions to generate HTML.

## Alternatives

There are a [number of template engines](https://awesome-go.com/#template-engines) available should you find `html/template` lacking for your use case. @SlinSo has put together a good [benchmark of these different template engines](https://github.com/SlinSo/goTemplateBenchmark).
