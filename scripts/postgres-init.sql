-- PostgreSQL Initialization Script for Temporal
-- =============================================

-- Create additional databases if needed
-- CREATE DATABASE afrikpay_analytics;

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create custom schemas for organization
CREATE SCHEMA IF NOT EXISTS analytics;
CREATE SCHEMA IF NOT EXISTS audit;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE temporal TO temporal;
GRANT ALL PRIVILEGES ON SCHEMA public TO temporal;
GRANT ALL PRIVILEGES ON SCHEMA analytics TO temporal;
GRANT ALL PRIVILEGES ON SCHEMA audit TO temporal;

-- Create audit table for transaction tracking
CREATE TABLE IF NOT EXISTS audit.transaction_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id VARCHAR(255) NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    user_id VARCHAR(255),
    amount DECIMAL(18,8),
    currency VARCHAR(10),
    status VARCHAR(50),
    external_reference VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for audit table
CREATE INDEX IF NOT EXISTS idx_transaction_logs_workflow_id ON audit.transaction_logs(workflow_id);
CREATE INDEX IF NOT EXISTS idx_transaction_logs_user_id ON audit.transaction_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_transaction_logs_status ON audit.transaction_logs(status);
CREATE INDEX IF NOT EXISTS idx_transaction_logs_created_at ON audit.transaction_logs(created_at);

-- Create analytics views
CREATE OR REPLACE VIEW analytics.daily_transactions AS
SELECT 
    DATE(created_at) as transaction_date,
    transaction_type,
    currency,
    COUNT(*) as transaction_count,
    SUM(amount) as total_amount,
    AVG(amount) as avg_amount
FROM audit.transaction_logs
WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE(created_at), transaction_type, currency
ORDER BY transaction_date DESC;

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for audit table
CREATE TRIGGER update_transaction_logs_updated_at 
    BEFORE UPDATE ON audit.transaction_logs 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample audit data for development
INSERT INTO audit.transaction_logs (
    workflow_id, 
    transaction_type, 
    user_id, 
    amount, 
    currency, 
    status,
    external_reference
) VALUES 
    ('wf_001', 'crypto_purchase', 'user_123', 100.00, 'USDT', 'completed', 'binance_ref_001'),
    ('wf_002', 'wallet_deposit', 'user_456', 50000.00, 'XAF', 'completed', 'mtn_ref_001'),
    ('wf_003', 'crypto_purchase', 'user_789', 0.001, 'BTC', 'pending', 'coinbase_ref_001');

-- Create performance monitoring view
CREATE OR REPLACE VIEW analytics.workflow_performance AS
SELECT 
    transaction_type,
    status,
    COUNT(*) as count,
    AVG(EXTRACT(EPOCH FROM (updated_at - created_at))) as avg_duration_seconds
FROM audit.transaction_logs
WHERE created_at >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY transaction_type, status
ORDER BY transaction_type, status;

-- Grant permissions on new objects
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA audit TO temporal;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA analytics TO temporal;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA audit TO temporal;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA analytics TO temporal;

-- Print completion message
DO $$
BEGIN
    RAISE NOTICE 'PostgreSQL initialization completed successfully!';
    RAISE NOTICE 'Schemas created: analytics, audit';
    RAISE NOTICE 'Tables created: audit.transaction_logs';
    RAISE NOTICE 'Views created: analytics.daily_transactions, analytics.workflow_performance';
    RAISE NOTICE 'Sample data inserted for development';
END $$;
