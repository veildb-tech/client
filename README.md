#DB Manager

This CLI application is used to download database dumps


Documentation:

---Creating a build for Linux
    To create a build for Linux, you need to go to the root of the directory where the code is contained and run the command: ```go build -o db-manager```
    And then place the executable file in the bin directory (if it doesn’t exist, then create it)

---Creating a build for Mac
    To create a build for Mac, you need to go to the root of the directory where the code is contained and run the command: ``` GOOS=darwin GOARCH=amd64 go build -o bin/db-mager-mac```
    ---ARM
        ```GOOS=darwin GOARCH=arm64 go build -o bin/db-mager-mac-arm```
    And then place the executable file in the bin directory (if it doesn’t exist, then create it)     

---Creating a build for Windows
    To create a build for Windows, you need to go to the root of the directory where the code is contained and run the command: ```GOOS=windows GOARCH=amd64 go build -o bin/db-manager-64.exe```

