# Rationales (crazy perhaps?)

If nothing else, please remember that influx is just a crazy experiment at this point.

The state of the art javascript ui frameworks have been born in one of the most actively hostile ecosystems in programming. There are more configurations and platform variations to deal with than any other environment given the state of browsers, javascript vms, javascript the language, and the dom. And yet people are developing projects like redux that distill and filter the experience of "developing an interactive web application" back to a sane situation.  Not taking inspirations from their creativity and ingenuity in developing libraries, development tools and especially debugging tools would be a mistake:  they've written some great code.

But go isn't javascript, and building redux for go wouldn't work well.  IMO, redux works because it feels well aligned with the needs for building a good highly interactive javascript application.   With that in mind, I trying to take the inspiration from redux's tooling awesomeness but build a framework that feels aligned with the needs for building a good, highly interactive go application.

---

The original goal is to build a fantastic development and debugging experience for gopherjs-based web front ends.  The present goal is to build a fantastic development and debugging experience for highly interactive applications.