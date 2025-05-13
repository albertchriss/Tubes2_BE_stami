<h1 align="center">STAMI</h1>
<h2 align="center">IF2211 Strategi Algoritma</h2>

## Table of Contents
1. [General Information](#general-information)
2. [Features](#features)
3. [Requirements](#requirements)
4. [Installation](#installation)
5. [Usage](#usage)
6. [Contributors](#contributors)


## General Information
Stami is a website that provides recipes for elements in the game Little Alchemy 2. It visualizes the recipe as a graph and uses graph search algorithms such as Breadth First Search (BFS), Depth First Search (DFS), and Bidirectional Search to find the solution. The website is built with Golang for the backend and Next.js for the frontend. Please note, that the forming elements are in lower tier than the initial element.
<br/><br/>
You can access our website in this link [stami](https://stami.cpoandy.me/home)!


## Features
| No. | Feature | Description |
|-----|----------------------|------------------------------------------------------------------|
| 1 | BFS | Systematically explores and visits all the vertices of a graph |
| 2 | DFS | Explores each branch as deeply as possible before backtracking |
| 3 | Bidirectional | Explores the search space from both the initial and goal nodes simultaneously |
| 4 | Multiple Recipe | Option to get multiple recipes for one element, if they exist |
| 5 | Live Update | View the search algorithm's process in real time |


## Requirements
Before setting up the project, make sure you have the following installed:
| No. | Required Program | Uses
| :--: | :--: | :--: |
| 1 | Makefile | Build automation tool |
| 2 | Docker | Containerization platform |
| 3 | Docker Compose | Tool for defining and running multi-container Docker application |
| 4 | Golang | Programming language used for the project |
| swag | Swagger Documentation generator for Go. [Link](https://github.com/swaggo/swag) |


### Installation

1. Clone the repository

    ```bash
    git clone https://github.com/albertchriss/Tubes2_BE_stami.git
    ```
2. Change directory to the project folder

    ```bash
    cd Tubes2_BE_stami
    ```
3. Create a `.env.dev` file in the `deployments/docker` directory and set the environment variables. You can use the `.env.example` file as a reference.

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

## Contributors
### **Kelompok Stami**
| NIM | Nama |
| :--: | :--: |
| 13523026 | Bertha Soliany Frandi |
| 13523077 | Albertus Christian Poandy |
| 13523089 | Ahmad Ibrahim |

Template by [Farhan Nabil Suryono](https://github.com/Altair1618)
