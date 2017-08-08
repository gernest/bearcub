[![Build Status](https://travis-ci.org/gernest/bearcub.svg?branch=master)](https://travis-ci.org/gernest/bearcub) [![Coverage Status](https://coveralls.io/repos/github/gernest/bearcub/badge.svg?branch=master)](https://coveralls.io/github/gernest/bearcub?branch=master) [![GoDoc](https://godoc.org/github.com/gernest/bearcub?status.svg)](https://godoc.org/github.com/gernest/bearcub)

**Background**:
  - https://help.stoplight.io/scenarios/overview/using-environments-and-variables

  The link above is a good example of what you will be implementing (or re-implementing). If you are familiar with html template systems, then think of this as a http request template system.

**Objective**:

  Coding is clearly important for this position, but there is so much more to this position than writing code. Just as important is your approach to solving the problem at hand - communication, code readability/documentation, and creativity are all important. We are excited to review your code, but are even more excited to understand the process you go through to solve the problem(s) at hand.

  1. Add the missing tests I described in variable_test.go.
  2. Implement variable replacement function.
  3. Benchmark variable replacement function.

**Requirements**:

  - Must use github.
  - No dependencies beyond the below - can only use golang's standard library.
    - You can use https://github.com/stretchr/testify to make testing easier.
    - You can use https://github.com/pytlesk4/m to help with JSON path selectors (For example, getting "users[0].id" in a JSON object).
  - Document your code. Bonus if you use godoc.
  - Must use go tools. I don't want to install any software.
  - Optimize variable replacement.
  - There should be no race conditions.
  - Should fail gracefully, should not panic.
  - Have fun! :)

**Bonus**:

  - Add variable_example_test.go.
  - Design/Propose an idea where you can use this (Write it down, draw it, peusdo code it, ...)
  - Do it.

**Instructions**:
  - Create a github repo, upload the attached folder, and invite me to it: https://github.com/pytlesk4
  - If you have any questions, please don't hesitate to ask. Github issues are a perfect way to communicate. So is email. So is slack (http://slack.stoplight.io), my username is pytlesk4, just direct message me with any questions. Or do a combination of the options above.
  - (Optional) Use github issues to break down the problem.
  - Use the standard libraries, no open source packages.
  - Before implementing the actual function set up all test cases, add any extra that I might have missed.
  - Implement variable replacement function.
  - Before you start actually optimizing your code, set up some benchmarks in variable_bench_test.go and send me the results. For this you can use whatever tools you want: benchcmp, pprof, uber's go-torch, ect... You aren't restricted to just go's benchmarking tool (it should be enough though). I want to run the benchmarks, so please include any instructions if you use something outside of go's benchmarking. Please also include the command you used to run it.
  - Optimize. That said, I don't expect you to get memory allocations down to zero, nor do I expect don't expect you to maximize every CPU Cycle. The more details you provide the better.
  - Save your final benchmarking results, and explain what you optimized and why.
