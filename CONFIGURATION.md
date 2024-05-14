# How to configure the agent

In the `config` directory, there is a `config.json.sample` file. This file contains an example of the configuration file that the agent uses. Rename this file to `config.json` and update the values to match your environment.

## Parameters

```json
{
    "server": {
        "port": ":9180", > The port that the agent will listen on.
        "LogRetentionDays": 2 > The number of days that the agent will keep the logs. The agent will delete logs older than this value.
    },
    "triggers": [
        {
            "Name":         "App1", > The application name, it is an arbitrary name that you can use to identify the project sinchronized via git. 
            "AbsolutePath": "/var/www/other/ap1", > The absolute path to the directory where the project is located. Best way to ge the right path is to use the `pwd` command in the terminal.
            "SharedSecret": "asjhhdal", > A arbitrary secret that you can use to secure the trigger endpoint. You will need to send this secret in the body of the request to the trigger endpoint.
            "SyncBranch":   "current", > current | specific  > The branch to be synchronized. If `current` is used, the agent will pull the current branch. If `specific` is used, the agent will check if the current brranch in the AbsolutePath matches the BranchName. If it does not match, the agent will return an error.
            "BranchName":   "main" > If `SyncBranch` is set to `specific`, the branch that should be selected on the server. 
        },
        {
            "Name":         "GA-DNS-PROXY",
            "AbsolutePath": "/var/www/app2/devops",
            "SharedSecret": "asjhhdal",
            "SyncBranch":   "specific", 
            "BranchName":   "main"
        }
    ]    
}
```

## Run the agent 

To start the agent on the server, go to the directory where the agent is located and run the following command:
`./sa_agent` or `./sa_agent &` to run the agent in the background.

## Stop the agent 

To stop the agent, you can use the `kill` command. First, you need to find the process ID of the agent. You can do this by running the following command:
`ps -ef | grep sa_agent`
then 
`kill <process_id>`

## Configuration changes

If the agent is running and you make changes to the configuration file, you will need to restart the agent for the changes to take effect. Just follow the steps above to stop the agent and then start it again.

# Usage 

The agent has 3 endpoints:
- `GET http://<your_public_ip>:<config.server.port>/info` - Returns the agent information
- `POST http://<your_public_ip>:<config.server.port>/trigger_update` - The endpoint to call to trigger the agent to run the configured operations. More on this Endpoint below.
- `ANY http://<your_public_ip>:<config.server.port>/` - Returns a 404 error for any other request.

## Trigger Update

To be able to trigger a `git pull orgin <branch>` on the configured directories, you need to send a POST request to the `trigger_update` endpoint. The request should contain a body with the following structure:

```json
{
    "Name": "App1",
    "SharedSecret": "asjhhdal"
}
```

The `Name` field should match the `Name` field in the `triggers` array in the configuration file. The `SharedSecret` field should match the `SharedSecret` field in the `triggers` array in the configuration file. If the `Name` and `SharedSecret` fields do not match the configuration file, the agent will return a 401 error.

### Error handling

If the agent encounters an error during the trigger operation, it will try to return an appropriate error message. If the agent is unable to determine the error, it will return a 500 error. For detialed error messages, check the agent logs on the server.