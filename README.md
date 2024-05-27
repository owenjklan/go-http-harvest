## What is this repository about?
As of May 2024, I have rekindled efforts to learn Golang.

This is an effort to make a Golang version of a Python exploration
tool that I wrote many years ago, "HTTP Harvest", which essentially
performs a rather dumb search for responsive HTTP servers on IP
addresses within a "Class C" network block. The Python version is
available [here, on Github](https://github.com/owenjklan/http_harvest) but
hasn't had much work done on it in years. It is functional but still has
signs of ways I was thinking of extending it. (I really need to clean those
up).

### The Pluralsight courses I've been using
I get no benefit from sharing this: Not in kahootz with Pluralsight in
any way. The two main courses that I've found to be of very high quality
were produced by author, Mike Van Sickle. In particular the courses are:
1. [Go Fundamentals](https://app.pluralsight.com/library/courses/fundamentals-go/table-of-contents)
2. [Concurrent Programming in Go](https://app.pluralsight.com/library/courses/go-programming-concurrent/table-of-contents)

According to Pluralsight, these courses were updated in November 2022 and February
2023, respectively.

### The main aspects this demonstrates
- Using the `flags` package for handling command line options
- Using "submodules" (still not 100% on correct terminologies). In particular the
  `byteseq` and `netobjects` directories.
- Unit testing Go code, as demonstrated in the "submodules" mentioned in the previous
  point.
- Basic use of the standard library's `http.HttpClient`
- Dispatching to specific functions based on a string-based map key. See `extras.go`
- Simple usage of both "fan-out" (single-producer, multiple-consumer) and "fan-in"
  (multiple-producer, single-consumer) concurrency patterns
- Manipulations and conversions of strings and IPv4 addresses (using `netip.Addr`
  instead of `net.IP` due to the former's ability to be comparable)

### Things that could definitely be improved
- There are a few Unit test cases that could be improved by using table-driven
  design and thus test move permutations than they currently do.
- The ability to output results in various file formats (plain, JSON, XML etc).
- The ability to add a scan delay interval so that requests don't smash the
  target network and raise suspicion (In this tools use case, the "slowness" of
  synchronous processing of the HTTP requests can actually be a benefit)
- Ability to specify that results should be stored to Amazon S3 or an SQS queue, because why not?
- Turn it into a simple webservice
- After making it a webservice, Dockerize it and then turn it into a Kubernetes
  Service / Deployment (for the purpose of more practical Kubernetes experience).
- The `byteseq` submodule is usable elsewhere, so it should probably be split into
  a module and repo of its own.