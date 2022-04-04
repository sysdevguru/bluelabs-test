# Bluelabs test
Simple Golang backend for wallet management.

## Assumption
- Don't need to take care of the transactions/users, just concentrate on wallet.
- Don't need to take care of the different currencies.
- User will have only one wallet.
- Don't need to take care of the authentication.

## Architectural decision
- Hexagonal architecture  
I used hexagonal architecture which is framework, library agnostic.  
Each `usecase` will have its own db repository, metrics etc and can be developed concurrently.

## How to run
```sh
sudo docker-compose up --remove-orphans
```
Listens on port `8080`

## Lint
```sh
sudo make tools
sudo make lint
```

## How to test
```sh
sudo make test
```

## Future improvements
- Create job queue for each tasks, `create_wallet`, `deposit`, `withdraw` and `get_balance`.
- Mock `wallet/Repo` interface
- Add function to get wallets of all users
- Add function to delete wallet
