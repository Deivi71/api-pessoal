-- Adicionar coluna foto_perfil na tabela usuarios
ALTER TABLE usuarios ADD COLUMN foto_perfil VARCHAR(255) DEFAULT NULL;

-- Comentário explicativo
-- Esta coluna armazenará o caminho/nome do arquivo da foto de perfil do usuário
-- Pode ser NULL se o usuário não tiver foto de perfil 