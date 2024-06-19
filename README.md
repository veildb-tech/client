# DBVisor 

DBvisor App is a CLI utility designed for convenient work on downloading the necessary database dumps
 
## Documentation:
### Flags
Flag ```-ldflags="-s -w"``` - is reduce the resulting binary size.
Flag ```-tags dev``` - served for compiling a binary file with settings for the dev server.

  
### Builds

#### Creating a build for Linux
To create a build for Linux, you need to go to the root of the directory where the code is contained and run the command: ```go build -ldflags="-s -w" -o dbvisor```
And then place the executable file in the bin directory (if it doesn’t exist, then create it)

#### Creating a build for Mac
To create a build for Mac, you need to go to the root of the directory where the code is contained and run the command: ```GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/dbvisor```

#### ARM
```GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/dbvisor```
And then place the executable file in the bin directory (if it doesn’t exist, then create it)

#### Creating a build for Windows
To create a build for Windows, you need to go to the root of the directory where the code is contained and run the command: ```GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/dbvisor.exe```

  
### Commands

- **login** - *Creating/updating a token and creating/editing a public key in the configuration file required for downloading database dumps.*
- **download** - *Downloading a dump of the database.*
- **save-key** - *Creating/editing a PEM public key.*
- **install** - *Installing a CLI application.*


