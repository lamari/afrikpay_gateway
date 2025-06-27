# Schémas de données CRUD

## User
```jsonc
{
  "id": "uuid",
  "email": "user@example.com",
  "role": "user", // enum: admin, user, operator, viewer
  "password_hash": "bcrypt$..."
}
```

## Wallet
```jsonc
{
  "id": "uuid",
  "user_id": "uuid",
  "balance": 0,
  "currency": "USDT" // enum: USDT, BTC
}
```

## Transaction
```jsonc
{
  "id": "uuid",
  "wallet_id": "uuid",
  "amount": 100,
  "type": "deposit", // enum: deposit, withdraw
  "status": "pending" // enum: pending, completed, failed
}
```
