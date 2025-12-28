-- TingHook MVP Initial Schema
-- See /plans/tinghook-mvp.md Section 7.2 Data Model

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL,
    subscription_plan VARCHAR(50) DEFAULT 'free',
    credits INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_api_key ON users(api_key);

-- Devices table
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    device_uid VARCHAR(64) UNIQUE NOT NULL,
    fcm_token TEXT,
    status VARCHAR(20) DEFAULT 'offline',
    battery_level INTEGER DEFAULT 0,
    app_version VARCHAR(20),
    last_seen_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_devices_user_id ON devices(user_id);
CREATE INDEX idx_devices_device_uid ON devices(device_uid);

-- Forwarding rules table
CREATE TABLE forwarding_rules (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    trigger_type VARCHAR(20) NOT NULL,
    sender_filter VARCHAR(255),
    content_filter TEXT,
    webhook_url TEXT NOT NULL,
    secret_header TEXT,
    method VARCHAR(10) DEFAULT 'POST',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_forwarding_rules_user_id ON forwarding_rules(user_id);
CREATE INDEX idx_forwarding_rules_device_id ON forwarding_rules(device_id);

-- Message logs table
CREATE TABLE message_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    direction VARCHAR(10) NOT NULL,
    sim_slot INTEGER DEFAULT 0,
    sender VARCHAR(50),
    receiver VARCHAR(50),
    content TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP
);
CREATE INDEX idx_message_logs_user_id ON message_logs(user_id);
CREATE INDEX idx_message_logs_created_at ON message_logs(created_at);
