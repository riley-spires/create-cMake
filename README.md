# CMake Project Generator

This Go program generates a CMakeLists.txt file for C++ projects and sets up the initial CMake configuration.

#### Features

- **Automatic CMakeLists.txt Generation**: Creates a CMakeLists.txt file based on user input and source files.
- **Customizable Project Settings**: Allows specifying project name, CMake version, and C++ standard.
- **Source File Validation**: Checks if specified source files exist before including them in the CMakeLists.txt.
- **CMake Configuration**: Attempts to build the CMake configuration after generating the CMakeLists.txt file.

#### Usage

```bash
go run main.go --project-name <name> {--cmake-version <version>} {--cxx-version <version>} <source_files...>
```

**Required Arguments:**
- `--project-name <name>`: Specifies the name of the project.

**Optional Arguments:**
- `--cmake-version <version>`: Sets the CMake version (default: 3.29).
- `--cxx-version <version>`: Sets the C++ standard version (default: 11).

**Source Files:**
- Provide one or more source file names as arguments.

#### Example

```bash
go run main.go --project-name MyProject --cmake-version 3.28 --cxx-version 17 main.cpp utils.cpp
```

This command will generate a CMakeLists.txt file for a project named "MyProject" using CMake version 3.28 and C++17 standard, including `main.cpp` and `utils.cpp` as source files.

#### Supported Platforms

Currently, the program supports automatic CMake configuration building on:
- Linux
- Windows

For other platforms, please submit a GitHub issue or pull request with your implementation.

#### Notes

- If a CMakeLists.txt file already exists, the program will prompt for confirmation before overwriting.
- The program validates the existence of specified source files before including them in the CMakeLists.txt.
- After generating the CMakeLists.txt, the program attempts to build the CMake configuration.

#### Building the Project

After successful CMake configuration, build the project using:

```bash
cmake --build build
```

#### Contributing

Contributions to add support for additional platforms or improve functionality are welcome. Please submit a pull request or create an issue for any enhancements or bug fixes.
