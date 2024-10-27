# ia04 - Decision-Making Solution Using Voting Algorithms
## Project Overview
### Description

This project, ia04, provides a decision-making solution implemented with various voting algorithms. The solution enables users to perform decision-making processes based on different social choice functions. Developed in Golang, this project includes several voting mechanisms, such as the Borda count, simple majority, Condorcet, and Copeland methods, which are accessible through a server API. The project demonstrates principles of agent-based modeling and social choice theory and allows for further experimentation with voting schemes.  

### Objectives

- Enable reliable, automated decision-making processes using voting algorithms.
- Provide modular and extendable code for experimenting with various social choice functions.
- Offer an accessible API for testing and utilizing voting systems.

### Features

- **Voting Algorithms**: Implements voting methods, including Borda, simple majority, Condorcet, Copeland, Approval...
- **Agent Interface**: Includes agent-based components for interfacing with the voting mechanisms.
- **API**: A server with endpoints to initiate voting processes and retrieve results.
- **Technologies Used**:  
    - Go (Golang) for programming
    - Go modules (go.mod and go.sum) for dependency management

## Project Structure
### Key Components

1. Command-line Tools (`cmd` directory):
    - `launch_server.go`: Launches the voting server.
    - `test_vote.go`, `test1_vote.go`: Command-line programs for testing the voting functionalities.

2. Voting Mechanisms (`comsoc` directory):
    - `approval.go`: Approval voting implementation.
    - `bordaMajority.go`: Borda voting implementation.
    - `copeland.go`: Copeland voting method.
    - `simpleMajority.go`: Simple majority voting.
    - `condorcet.go`: Condorcet method implementation.
    - `utils.go`, types.go: Utility functions and type definitions for voting mechanisms.
    - `tiebreak.go`: Handles tie-breaking procedures.

3. Agent Components (`agt` directory):
    - `agentInterface.go`: Defines interfaces for agents.
    - `agentImplementation.go`: Implementation of the agent interface.
    - `types.go`: Types associated with agent interactions.

4. Server Components (`server` directory):
    - `result.go`, `vote.go`, `ballot.go`: Core functionalities related to ballots and votes.
    - `handlers.go`: API endpoint handlers.
    - `routes.go`: Route definitions for the server.

## Implementation
### Project Files

- `go.mod`, `go.sum`: Dependency management files.
- Go source files for implementing voting mechanisms, agents, and server endpoints.
- Command-line tool scripts for testing and launching the server.

### Dependencies

- Golang: Required for building and running the project.

### Usage

Set Up the Environment:
- Ensure Go is installed and set up on your system.

Clone the repository and navigate to the project directory :

```bash
git clone git@github.com:TobiasInfo/SystemeMultiAgents.git
cd SystemeMultiAgents/ia04
```

Run the Server:

```bash
go run cmd/launch_server/launch_server.go
```
Test Voting Algorithms:

Run individual tests or experiments with:

```bash
go run cmd/test_vote/test_vote.go
go run cmd/test1_vote/test1_vote.go
```

## Learning Outcomes
### Technical Skills

- Implementation of social choice functions
- API development with Golang
- Agent-based design patterns

### Soft Skills

- Problem-solving in algorithm design
- Collaboration on modular software development

## Challenges and Improvements
### Challenges

- Understanding the intricacies of social choice theory.
- Implementing tie-breaking mechanisms across different voting methods.

### Improvements

- Enhance algorithm efficiency.
- Extend the server with additional endpoints for further voting scenarios.

## Authors

- Tobias SAVARY
- Nassim SAIDI

## License

This project is licensed under the MIT License - see the LICENSE file for details.