# klaw

![klaw logo](https://github.com/joaooliveirapro/klaw/blob/master/klaw.png)
**klaw** is a command-line tool written in Go that scans your code for TODOs and integrates directly with GitHub to create issues. 

It's fast, lightweight `(7Mb)`, and designed to help developers keep track of tasks in their projects.

Inspired by [@tsoding/snitch](https://github.com/tsoding/snitch)

---

## Features

- Scan directories for code comments like `// TODO: ...` 
- Lists all open and closed issues
- Create GitHub issues automatically from untracked TODOs
- Colored terminal output with ASCII logo  
- Cross-platform executable with custom icon  
- Easy CLI usage with flags 

---

## Installation
### From source

```bash
git clone https://github.com/joaooliveirapro/klaw.git
cd klaw
go build -ldflags "-s -w" -o klaw.exe

# usage
klaw -l <directory_path> <file_extension> 
```
