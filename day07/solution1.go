package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
)

type File struct {
    Name string
    Size int
}

func NewFile(name string, size int) *File {
    return &File{name, size}
}



type Dir struct {
    Name string
    Dirs []*Dir
    Files []*File
    Parent *Dir
}

func NewDir(name string) *Dir {
    return &Dir{
        name,
        make([]*Dir, 0),
        make([]*File, 0),
        nil,
    }
}

func (d *Dir) AddDir(newDir *Dir) {
    newDir.Parent = d
    d.Dirs = append(d.Dirs, newDir)
}

func (d *Dir) AddFile(newFile *File) {
    d.Files = append(d.Files, newFile)
}

func (d *Dir) ExistsFile(name string) bool {
    for _, file := range d.Files {
        if file.Name == name {
            return true
        }
    }
    return false
}

func (d *Dir) ExistsDir(name string) bool {
    for _, dir := range d.Dirs {
        if dir.Name == name {
            return true
        }
    }
    return false
}

func (d *Dir) GetDir(name string) *Dir {
    for _, dir := range d.Dirs {
        if dir.Name == name {
            return dir
        }
    }
    return nil
}

func (d *Dir) Size() int {
    var total = 0
    for _, file := range d.Files {
        total += file.Size
    }
    for _, dir := range d.Dirs {
        total += dir.Size()
    }
    return total
}

func (d *Dir) Walk(observe func(*Dir)) {
    observe(d)
    for _, dir := range d.Dirs {
        dir.Walk(observe)
    }
}

func EachLineDo(f func(string)) error {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        line := scanner.Text()
        f(line)
    }

    err := scanner.Err()

    return err
}

func main() {

    root := NewDir("/")
    cwd := root
    lastCommand := ""

    err := EachLineDo(func(line string) {
        if line[0] == '$' {
            split := strings.Split(line, " ")
            cmd := split[1]
            args := strings.Join(split[2:], " ")

            switch cmd {
            case "cd":
                switch args {
                case "/":
                    cwd = root
                case "..":
                    cwd = cwd.Parent
                default:
                    dir := cwd.GetDir(args)
                    if dir == nil {
                        dir = NewDir(args)
                        cwd.AddDir(dir)
                    }
                    cwd = dir
                }
            case "ls":
                ;
            default:
                ;
            }

            lastCommand = cmd
        } else {
            if lastCommand == "ls" {
                split := strings.Split(line, " ")
                if split[0] == "dir" {
                    dirName := split[1]
                    if !cwd.ExistsDir(dirName) {
                        cwd.AddDir(NewDir(dirName))
                    }
                } else {
                    fileName := split[1]
                    if !cwd.ExistsFile(fileName) {
                        size, err := strconv.Atoi(split[0])
                        if err != nil {
                            fmt.Fprintf(os.Stderr, "not an integer", split[0])
                            os.Exit(1)
                        }
                        cwd.AddFile(NewFile(fileName, size))
                    }
                }
            }
        }
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "reading stdin:", err)
        os.Exit(1)
    }

    var total = 0
    root.Walk(func(cwd *Dir) {
        var size = cwd.Size()
        if size <= 100000 {
            total += size
        }
    })

    fmt.Println("Total size:", total)
}
