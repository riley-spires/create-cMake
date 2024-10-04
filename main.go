package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"os/exec"
	"runtime"
	"bufio"
)

func main() {
	if _, err := os.Stat("CMakeLists.txt"); err == nil {
		scanner := bufio.NewScanner(os.Stdin)

		valid_options := []string{"yes", "no", "y", "n"}

		for !slices.Contains(valid_options, strings.ToLower(scanner.Text())) {
			fmt.Println("CMakeLists.txt already exists")
			fmt.Print("Are you sure you want to override (yes|no): ")
			scanner.Scan()
		}

		if strings.ToLower(scanner.Text()) == "n" || strings.ToLower(scanner.Text()) == "no" {
			fmt.Fprintln(os.Stderr, "Aborting...")
			return
		}
	}

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

	if len(flag.Args()) == 0 {
		fmt.Fprint(os.Stderr, "ERROR: Must provide at least one source file!\n")
		return
	}

	for _, arg := range flag.Args() {
		if _, err := os.Stat(arg); errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "ERROR: source file %s does not exist!\n", arg)
			return
		}
	}

	fmt.Println("Attempting to create CMakeLists.txt with following config:")
	fmt.Printf("Project Name: %s\n", project_name)
	fmt.Printf("CMake Version: %s\n", cmake_version)
	fmt.Printf("Cxx Version: %s\n", cxx_version)
	fmt.Println("Source files:")
	for _, arg := range flag.Args() {
		fmt.Printf("\t%s\n", arg)
	}

	file, err := os.Create("./CMakeLists.txt")

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not create CMakeLists.txt : %s\n", err.Error())
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

	fmt.Println("Created CMakeLists.txt")
	
	os_name := runtime.GOOS

	fmt.Printf("%s OS detected. Attempting to build cmake config\n", os_name)

	var cmd *exec.Cmd = nil

	switch os_name {
		case "linux":
			cmd = exec.Command("bash", "-c", "cmake . -B build")
			break
	}

	if cmd == nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s is an unsupported platform!\n", os_name)
		fmt.Fprint(os.Stderr, "Please make a github issue or pull request!\n")
		return
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not build cmake config : %s\n", err.Error())
		return
	}
	
	fmt.Println("Cmake config built. Build project with \"cmake --build build\"")
}
