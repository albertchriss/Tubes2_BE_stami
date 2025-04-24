# Tubes 2 STIMA
## Table of Contents

- [Tubes 2 Stima](#tubes-2-stima)
    - [Table of Contents](#table-of-contents)
    - [Getting Started](#getting-started)
        - [Prerequisites](#prerequisites)
        - [Installation](#installation)
    - [Usage](#usage)
    - [Acknowledgements](#acknowledgements)

## Getting Started

### Prerequisites

Before setting up the project, make sure you have the following installed:

- **Makefile** - Build automation tool.
- **Docker** - Containerization platform.
- **Docker Compose** - Tool for defining and running multi-container Docker applications.
- **Docker Desktop** (optional) - Docker GUI for managing containers.
- **Golang** - Programming language used for the project.
- **swag** - Swagger Documentation generator for Go. https://github.com/swaggo/swag


### Installation

1. Clone the repository

    ```bash
    git clone https://github.com/albertchriss/Tubes2_BE_stami.git
    ```
2. Change directory to the project folder

    ```bash
    cd Tubes2_BE_stami
    ```
3. Create a `.env.dev` file in the root directory and set the environment variables. You can use the `.env.example` file as a reference.

4. Install the required Go modules

    ```bash
    go mod download
    ```
    or
    
    ```bash
    go mod tidy
    ```

## Usage

### Run the Application
Build the Docker image and run the container
```bash
make up-dev
```
Remove the container
```bash
make down-dev
```
### Generate Swagger Documentation
```bash
make generate-docs
```

## Acknowledgements

This project is being developed by:

- Bertha Soliany Frandi - 13523026
- Albertus Christian Poandy - 13523077
- Ahmad Ibrahim - 13523089

Template by [Farhan Nabil Suryono](https://github.com/Altair1618)
