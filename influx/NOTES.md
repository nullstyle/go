# Rationales (crazy perhaps?)

If nothing else, please remember that influx is just a crazy experiment at this point.

The state of the art javascript ui frameworks have been born in one of the most actively hostile ecosystems in programming. There are more configurations and platform variations to deal with than any other environment given the state of browsers, javascript vms, javascript the language, and the dom. And yet people are developing projects like redux that distill and filter the experience of "developing an interactive web application" back to a sane situation.  Not taking inspirations from their creativity and ingenuity in developing libraries, development tools and especially debugging tools would be a mistake:  they've written some great code.

But go isn't javascript, and building redux for go wouldn't work well.  IMO, redux works because it feels well aligned with the needs for building a good highly interactive javascript application.   With that in mind, I trying to take the inspiration from redux's tooling awesomeness but build a framework that feels aligned with the needs for building a good, highly interactive go application.

---

The original goal is to build a fantastic development and debugging experience for gopherjs-based web front ends.  The present goal is to build a fantastic development and debugging experience for highly interactive applications.


# An Alternate take on the vision for this project (as an exercise to guide the design decisions of the package structure)

An alternate go standard library for developing highly interactive applications. In most senses when we say standard library we mean "the set of packages provided by the language developers to solve problems general to developers using the language", but for this project I am proposing a standard library to mean "a set of packages designed to optimize the experience of develop, operate and maintain a certain classification of applications". In this sense, I can base the design of my packages to mirror the go standard library to help build intuition about the frameworks architecture.  For example, any go programmer should know that the "fmt" package is used for formatting values into strings and know of the functions `Sprintf` and siblings.  Then, if we design our packages are of the same shape, they should be more likely to remember that the `influxfmt` package is used for formatting influx-specific values into strings.  (aside: hopefully this architecture makes it easier for others to contribute to  this project because a contribution could be as simple as "transliterating" an actual go stdlib function into my poorly covered versions of the various stdlib packages)

Said another way, I aim to build a framework used to build highly interactive applications that is so powerful it can convince a javascript developer to switch to  go and get good at it enough to leverage the similarities between my framework and the go standard library.  At the very least having this vision means my design will be heavily informed by people far smarter than me (the go team).


# Sharing through composition

Component developers should design their components to be configured by an ancestor component by exposing fields that can be configured at load time. Components should abstract complexity in the action stream by providing  methods that an ancestor can call during dispatch.