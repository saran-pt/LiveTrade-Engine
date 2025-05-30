# LiveTrade Engine – A Real-Time Trading Engine in Go

## Objective

**LiveTrade Engine** is a real-time order matching engine developed in Golang, designed to simulate trading environments like stock exchanges or crypto order books. It processes buy/sell orders, matches them based on defined rules (e.g. price-time priority), and maintains a live order book.

## Data Handling

User-submitted orders act as the live data for this engine. Each order is processed in real-time either matched immediately with an existing order or stored in the order book for future matching. Successfully matched trades and active orders are tracked to reflect the current state of the trading system.

## Architecture & Workflow

1. Users submit orders via REST API.
2. The system accepts the order, validates it, and adds it to the order book (in-memory).
3. The matching engine attempts to find a suitable match based on price and side.
4. Matched trades are recorded in the database; unmatched orders stay in memory.
5. APIs are provided to fetch current order book state.

## Future Enhancements

- [ ] Add WebSocket support for live order updates
- [ ] Implement user authentication
- [ ] Introduce more complex matching rules (market/limit orders, partial fills)
- [ ] Add dashboard to visualize trading flow
