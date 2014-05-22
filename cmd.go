package main

import (
        "os/exec"
        "strconv"
)

func Configure() error {
        cmd := exec.Command("sh", "./nginx-configure")
        CheckVerboseEnabled(cmd)
        return cmd.Run()
}

func Make(jobs int) error {
	cmd := exec.Command("make", "-j", strconv.Itoa(jobs))
        CheckVerboseEnabled(cmd)
        return cmd.Run()
}

func ExtractArchive(path string) error {
        cmd := exec.Command("tar", "zxvf", path)
        CheckVerboseEnabled(cmd)
        return cmd.Run()
}


func SwitchRev(rev string) error {
        cmd := exec.Command("git", "co", rev)
        CheckVerboseEnabled(cmd)
        return cmd.Run()
}
