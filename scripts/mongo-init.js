// MongoDB Initialization Script for Afrikpay Gateway
// ================================================

// Switch to afrikpay database
db = db.getSiblingDB('afrikpay');

// Create collections with validation
db.createCollection('users', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['email', 'created_at'],
      properties: {
        email: {
          bsonType: 'string',
          pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        },
        phone: {
          bsonType: 'string'
        },
        created_at: {
          bsonType: 'date'
        },
        updated_at: {
          bsonType: 'date'
        }
      }
    }
  }
});

db.createCollection('wallets', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['user_id', 'currency', 'balance', 'created_at'],
      properties: {
        user_id: {
          bsonType: 'objectId'
        },
        currency: {
          bsonType: 'string',
          enum: ['USD', 'XAF', 'USDT', 'BTC']
        },
        balance: {
          bsonType: 'decimal'
        },
        created_at: {
          bsonType: 'date'
        },
        updated_at: {
          bsonType: 'date'
        }
      }
    }
  }
});

db.createCollection('transactions', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['user_id', 'type', 'amount', 'currency', 'status', 'created_at'],
      properties: {
        user_id: {
          bsonType: 'objectId'
        },
        type: {
          bsonType: 'string',
          enum: ['crypto_purchase', 'wallet_deposit', 'transfer']
        },
        amount: {
          bsonType: 'decimal'
        },
        currency: {
          bsonType: 'string',
          enum: ['USD', 'XAF', 'USDT', 'BTC']
        },
        status: {
          bsonType: 'string',
          enum: ['pending', 'processing', 'completed', 'failed', 'cancelled']
        },
        created_at: {
          bsonType: 'date'
        },
        updated_at: {
          bsonType: 'date'
        }
      }
    }
  }
});

// Create indexes for performance
db.users.createIndex({ 'email': 1 }, { unique: true });
db.users.createIndex({ 'phone': 1 }, { unique: true, sparse: true });
db.users.createIndex({ 'created_at': 1 });

db.wallets.createIndex({ 'user_id': 1, 'currency': 1 }, { unique: true });
db.wallets.createIndex({ 'user_id': 1 });

db.transactions.createIndex({ 'user_id': 1 });
db.transactions.createIndex({ 'status': 1 });
db.transactions.createIndex({ 'type': 1 });
db.transactions.createIndex({ 'created_at': 1 });

// Insert sample data for development
db.users.insertMany([
  {
    _id: ObjectId(),
    email: 'john.doe@example.com',
    phone: '+237123456789',
    first_name: 'John',
    last_name: 'Doe',
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    _id: ObjectId(),
    email: 'jane.smith@example.com',
    phone: '+237987654321',
    first_name: 'Jane',
    last_name: 'Smith',
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print('MongoDB initialization completed successfully!');
print('Collections created: users, wallets, transactions');
print('Indexes created for optimal performance');
print('Sample data inserted for development');
