CREATE INDEX idx_oauth_provider_external ON oauth_providers(provider_name, external_id);
CREATE INDEX idx_oauth_student ON oauth_providers(student_id);
