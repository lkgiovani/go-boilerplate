-- Migration: V4 - Criação da tabela de tokens de verificação de email
-- Descrição: Tabela para armazenar tokens de verificação de email dos usuários

CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP(6) NOT NULL,
    verified_at TIMESTAMP(6),
    used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_email_verification_tokens_user 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices para otimização
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_id ON email_verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_token ON email_verification_tokens(token);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_expires_at ON email_verification_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_email ON email_verification_tokens(email);
CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_used ON email_verification_tokens(used);

-- Comentários para documentação
COMMENT ON TABLE email_verification_tokens IS 'Tokens para verificação de email dos usuários';
COMMENT ON COLUMN email_verification_tokens.user_id IS 'Referência ao usuário';
COMMENT ON COLUMN email_verification_tokens.email IS 'Email que está sendo verificado';
COMMENT ON COLUMN email_verification_tokens.token IS 'Token único de verificação';
COMMENT ON COLUMN email_verification_tokens.expires_at IS 'Data e hora de expiração do token';
COMMENT ON COLUMN email_verification_tokens.verified_at IS 'Data e hora da verificação';
COMMENT ON COLUMN email_verification_tokens.used IS 'Indica se o token foi usado';
