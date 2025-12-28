-- Rollback TingHook MVP Initial Schema
-- Drop in reverse order to respect foreign key constraints

DROP INDEX IF EXISTS idx_message_logs_created_at;
DROP INDEX IF EXISTS idx_message_logs_user_id;
DROP TABLE IF EXISTS message_logs;

DROP INDEX IF EXISTS idx_forwarding_rules_device_id;
DROP INDEX IF EXISTS idx_forwarding_rules_user_id;
DROP TABLE IF EXISTS forwarding_rules;

DROP INDEX IF EXISTS idx_devices_device_uid;
DROP INDEX IF EXISTS idx_devices_user_id;
DROP TABLE IF EXISTS devices;

DROP INDEX IF EXISTS idx_users_api_key;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
