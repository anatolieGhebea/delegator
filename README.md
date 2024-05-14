# simple-automation-agent
@todo


# BUILD Project

## Linux / amd64
`GOOS=linux GOARCH=amd64 go build -o ./dist/linux/<version>/sa_agent`

## MacOS / 
This agent is intended to run on the server, which are usually linux based. However, if you want to run the agent on your local machine, you can build the agent for MacOS.
`GOOS=darwin GOARCH=arm64 go build -o ./build/darwin/version/sa_agent`

# Create compiled for distribution

## Linux / amd64
After compiling the project, in the same dist directory copy the `/config/config.json.sample` then use tar to zip the file for easier distribution.