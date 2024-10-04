package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"os/exec"
)

func main() {
	required_args := []string{"project-name"}
	project_name := ""
	cmake_version := "3.29"
	valid_cmake_versions := []string{"3.20", "3.21", "3.22", "3.23",
		"3.24", "3.25", "3.26", "3.27",
		"3.28", "3.29", "3.30",
	}
	cxx_version := "11"
	valid_cxx_versions := []string{"98", "03", "11", "14", "17", "20", "23", "26"}

	flag.Func("project-name", "-project-name {name}", func(flag_val string) error {
		if flag_val == "" {
			return errors.New("project-name cannot be blank")
		}

		if project_name != "" {
			return errors.New("project-name flag already provided")
		}

		project_name = flag_val

		return nil
	})

	flag.Func("cmake-version", "-cmake-version {3.20, 3.21, 3.22, 3.23, 3.24, 3.25, 3.26, 3.27, 3.28, 3.29, 3.30}", func(flag_val string) error {
		if flag_val == "" {
			return errors.New("cmake-version cannot be blank")
		}

		if !slices.Contains(valid_cmake_versions, flag_val) {
			return errors.New("unsupported cmake version")
		}

		cmake_version = flag_val

		return nil
	})

	flag.Func("cxx-version", "-cxx-version {98, 03, 11, 14, 17, 20, 23, 26}", func(flag_val string) error {
		if flag_val == "" {
			return errors.New("cxx-version cannot be blank")
		}

		if !slices.Contains(valid_cxx_versions, flag_val) {
			return errors.New("unsupported cxx version")
		}

		cxx_version = flag_val

		return nil
	})

	for _, arg := range required_args {
		found := false
		for _, cmd_arg := range os.Args {
			if strings.Contains(cmd_arg, arg) { found = true }
		}
		if !found {
			fmt.Fprintf(os.Stderr, "ERROR: must provide a value with --%s\n", arg)
			return
		}
	}

	flag.Parse()

	file, err := os.Create("./CMakeLists.txt")

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not create CMakeLists.txt -> %s\n", err.Error())
	}

	file.WriteString(fmt.Sprintf("cmake_minimum_required(VERSION %s)\n", cmake_version))
	file.WriteString(fmt.Sprintf("project(%s)\n\n", project_name))
	file.WriteString(fmt.Sprintf("set(CMAKE_CXX_STANDARD %s)\n", cxx_version))
	file.WriteString("set(CMAKE_CXX_STANDARD_REQUIRED True)\n\n")
	file.WriteString("set(SOURCES ")
	for _, file_name := range flag.Args() {
		file.WriteString(fmt.Sprintf("\"%s\"\n", file_name))
	}
	file.WriteString(")\n\n")
	file.WriteString("add_executable(${PROJECT_NAME} ${SOURCES})")
	file.Close()

	cmd := exec.Command("bash", "-c", "cmake . -B build")

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not build cmake config : %s\n", err.Error())
	}
}
