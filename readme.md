# klaw

<img src="https://github.com/joaooliveirapro/klaw/blob/master/klaw.png" width="128" height="128" />


**klaw** is a command-line tool written in Go that scans your code for `// TODO: ... ` and integrates directly with GitHub to create issues. 

It's lightweight `(7Mb)`, fast and designed to help developers keep track of tasks in their projects.

Inspired by [@tsoding/snitch](https://github.com/tsoding/snitch)

---

## Features

- Scan directories for code comments like `// TODO: ...` 
- Lists all open and closed issues
- Create GitHub issues automatically from untracked `// TODO: ... `
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
