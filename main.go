package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
    current_env := runtime.GOOS

    var command string
    var remotes_url string
    var err error

    // Check the current runtime so that we know which command to use to open the link in main browser
    switch current_env {
    case "macos":
        command = "open"
    case "window":
        fmt.Errorf("Can't run this on windows at the moment")
    default:
        command = "xdg-open"
    }

    // KDE and Gnome might use other commands for default application.
    // Wayland also might use something else
    // Need to test this and add logic for it if xdg-open is not available in linux always
    // Somtimes linux might not have xdg-open command forexample on wayland
    if command == "xdg-open" {
        cmd := exec.Command(command)
        if err = cmd.Start(); err != nil {
            fmt.Printf("%v", err)
        }
    }

    // Save the git remotes so we can manipulate them later to generate the link
    out, err := exec.Command("git", "ls-remote", "--get-url").Output()
    if err != nil {
        fmt.Printf("%v", err)
    }

    remotes := string(out[:])

    // TODO: Add option to open push url instead of fetch url
    // Some git origins are still https instead of ssh so we need to check
    // if the remote contains an @ sign or not
    if strings.Contains(remotes, "@") {
        remotes_url = strings.Split(strings.Split(remotes, "@")[1], " ")[0]
    } else {
        // Here we split on "://" so we can later have similar strings later
        remotes_url = strings.Split(strings.Split(remotes, "://")[1], " ")[0]
    }

    // ssh url contains : which needs to be changed slash so we do that here
    if strings.Contains(remotes_url, ":") {
        remotes_url = strings.Replace(remotes_url, ":", "/", 1)
    }

    // Save prefix and url to array of strings so we can use join on them later
    url_parts := []string{"https://", remotes_url}
    out, err = exec.Command(command, strings.Join(url_parts, "")).Output()

    if err != nil {
        fmt.Printf("%v", err)
    }
}

