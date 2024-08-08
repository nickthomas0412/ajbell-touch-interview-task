# AJBell Touch Interview Task
Interview task for the Software Engineer role at AJBell.

## Prerequisites
- Go (Tested with 1.22.5)
- Git

## Installation
1. Clone the repository
```
git clone https://github.com/nickthomas0412/ajbell-touch-interview-task.git
cd ajbell-touch-interview-task
```

2. Download required modules
```
go mod tidy
```

3. Run the backend server
```
cd src
go run main.go
```

4. Run the example script in a new terminal window
```
cd src
cd examples
go run example.go
```

5. Build the server
```
go build -o backend-server
```

6. Run the server
```
./backend-server
```
