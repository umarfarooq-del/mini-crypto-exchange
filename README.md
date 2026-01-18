# 📘 Limit Order Book Matching Engine (Go)

## Overview

This project implements a **limit order book matching engine** similar to a simplified crypto exchange.
It allows users to place **limit BUY and SELL orders** for a trading pair (e.g. `BTC/USDT`) and automatically matches them using **price–time priority (FIFO)**.

The focus of this implementation is **correctness of matching logic**, clean in-memory design, and clear separation of concerns.

---

## Features

- ✅ Limit BUY and SELL orders only
- ✅ Price priority (best price first)
- ✅ Time priority (FIFO for same price)
- ✅ Partial and full order matching
- ✅ In-memory order book
- ✅ In-memory order history
- ✅ Trade creation on match
- ✅ Thread-safe matching engine
- ✅ REST APIs

---

## API Endpoints

### 1️⃣ Create Trading Pair (Admin)

```
POST /api/pairs
```

**Request**
```json
{
  "base": "BTC",
  "quote": "USDT"
}
```

Creates a tradable pair (`BTC/USDT`) and initializes its order book.

---

### 2️⃣ Place Order

```
POST /api/orders
```

**Request**
```json
{
  "pair": "BTC/USDT",
  "side": "buy",
  "price": 13000,
  "quantity": 1,
  "user_id": 101
}
```

**Rules**
- Only **limit orders** are supported
- BUY → maximum price user is willing to pay
- SELL → minimum price user is willing to accept

**Response**
```json
{
  "order": { ... },
  "trades": [ ... ]
}
```

---

### 3️⃣ Get User Orders

```
GET /api/orders?user_id=101
```

Returns **all orders placed by the user**, including:
- `open`
- `partial`
- `filled`

---

### 4️⃣ Get Order Book

```
GET /api/orderbook?pair=BTC/USDT
```

Returns the current **order book snapshot**, showing only open and partially filled orders.

---

## Core Matching Logic

### Price Priority
- BUY orders match against the **lowest priced SELL**
- SELL orders match against the **highest priced BUY**

### Time Priority (FIFO)
- If prices are equal, the **earliest created order** is matched first

### Trade Price Rule
> Trades execute at the **price of the existing order in the order book**, not the incoming order.

---

## Internal Design

```
MatchingEngine
 ├── orderBooks (per trading pair)
 ├── orders     (all orders, in memory)
 ├── trades     (executed trades)
```

Filled orders are removed from the order book but retained in order history.

---

## Thread Safety

- Uses `sync.RWMutex`
- Write lock for order placement and matching
- Read lock for order history queries

---

## Storage Strategy

- Entire system is **in memory**
- State is lost on service restart
- Design allows database integration later with minimal changes

---

## Setup Instructions
go mod tidy

### Run the Server
```bash
go run cmd/server/main.go
```

Server starts on:
```
http://localhost:50053
```

---

## Summary

This project demonstrates:
- Correct limit order book matching
- Proper handling of partial fills
- Clean separation of concerns
- Thread-safe in-memory design
