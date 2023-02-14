# milton
Milton is a Slack bot written in Go that helps manage code review interrupts.  
It provides a set of commands that can be used within Slack and consumes Slack's RTM event stream by running in socket mode

## Slack Commands
Milton responds to the following Slack commands:

- `!queue` Adds items to an interrupt queue for the given Slack channel, example:
```
!queue https://github.com/pulls/1
!queue https://github.com/pulls/1 https://github.com/pulls/2
!queue https://github.com/pulls/3 // this is regarding feature A
```
- `!return` Returns items from the interrupt queue for the given Slack channel and notifies the requester.
Returned items means the reviewer saw some concerns and would like the requester to provide more context or update code.
Usage example
```
!return https://github.com/pulls/1
!return https://github.com/pulls/1 https://github.com/pulls/2
!return https://github.com/pulls/3 // this is returned, please look at feedback in PR
```
- `!done` Removes an item from the interrupt queue and sends a notification to the requester, signalling the reviewer is content with the code.
```
!done https://github.com/pulls/1
!done https://github.com/pulls/1 https://github.com/pulls/2
!done https://github.com/pulls/3 // looks good, nit in PR, reach out to team X for deploy etc.
```
- `!help` Displays a list of available commands.

## Installation
To install and use Milton, follow these steps:

- Clone this repository.
- Install the required dependencies using `go mod vendor`
- Create a `config.yaml` file in the root directory of the project with the following structure:
```yaml
database_metadata:
  type: mysql
  host: localhost
  name: milton
  user: milton
  password: milton
slack:
  app_token: xapp-token
  bot_token: xoxb-token
```
- Build the application using `go build cmd/milton/milton.go`

## Plugin Development
Milton is designed to be extensible, and developers can add their own Slack commands and workflows.  
To add a new plugin, create a new Go module under the `src/plugins` directory.

The plugin must implement the `slack.Command` interface and define its own set of Slack commands.  
The implementation should be placed in a separate file, and the module should be imported in `pkg/milton/milton.go`

To selectively load the plugin when running Milton,  
update the `setupPlugins` method in `pkg/milton/milton.go` to load the module appropriately.

You can also use existing plugins as examples for developing your own.  
For example, to see the implementation of the OpsGenie Slack command, see the [opsgenie slack command setup](src/plugins/opsgenie) and [opsgenie api usage](src/opsgenie).  

Once you have created the plugin, you can selectively have it be loaded when running Milton:
```bash
./milton run -x plugin-name -c config.yaml -d src/backend/models
```
Where `plugin-name` is the name of the module containing the plugin implementation.

## Usage
Milton can be used through the command-line interface. The following commands are available:

- `milton run` Starts the bot and listens for Slack events
To see a list of available flags and options for a command, run:

```bash
./milton [command] --help
```
For example, to see the available options for the run command, run:
```bash
./milton run --help
```

## Local Testing
### Build Docker Container
```
docker build -t milton .
```
### Bring Up Stack
```
docker-compose up -d
```

### Teardown Stack
```
docker-compose down
docker volume rm milton_db_data
```

## Contribution
Contributions to Milton are welcome! If you would like to contribute, please submit a pull request.
