# Erc20 Event Logger 
## Log ERC20 transfer events to a database

Simple Go service that:

Listens for ERC20 Transfer events on Ethereum or any EVM-compatible chain

Stores events in SQLite database

Provides one API endpoint to fetch all logged events

## How to Use
Clone repo:
git clone https://github.com/your-username/erc20-event-logger-go.git

Install dependencies:
go get

Set up:

Create .env file with your Infura WebSocket URL:
INFURA_WSS=wss://mainnet.infura.io/ws/v3/your-project-id

Run:
go run main.go
